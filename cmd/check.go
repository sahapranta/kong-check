package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/sahapranta/kong-check/config"
	"github.com/sahapranta/kong-check/db"
	"github.com/sahapranta/kong-check/models"
)

func CheckRoutes(conf *config.Config, customPath string, verbose bool, timeout int, checkAll bool, serviceNames string, methods string) {
	var routes []models.Route
	var err error

	if checkAll {
		routes, err = db.GetAllRoutes(conf)
		if err != nil {
			log.Fatalf("Failed to query routes: %v", err)
		}
	} else if serviceNames != "" {
		services := strings.Split(serviceNames, ",")
		routes, err = db.GetRoutesByServiceNames(conf, services)
		if err != nil {
			log.Fatalf("Failed to query routes for services %s: %v", serviceNames, err)
		}
	} else {
		fmt.Println("Please specify either --all to check all routes or --services to check specific services")
		return
	}

	if len(routes) == 0 {
		fmt.Println("No routes found to check.")
		return
	}

	fmt.Printf("Checking %d routes (timeout: %ds)...\n\n", len(routes), timeout)

	// Set up client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Define colors
	successColor := color.New(color.FgGreen, color.Bold).SprintFunc()
	errorColor := color.New(color.FgRed, color.Bold).SprintFunc()
	serviceColor := color.New(color.FgCyan).SprintFunc()
	// routeColor := color.New(color.FgYellow).SprintFunc()
	methodColor := color.New(color.FgMagenta).SprintFunc()
	timeColor := color.New(color.FgBlue).SprintFunc()

	sem := make(chan bool, conf.MaxConcurrentReq)
	var wg sync.WaitGroup

	// Channel for results
	results := make(chan models.CheckResult, len(routes)*2) // Buffer size to prevent blocking

	go func() {
		for result := range results {
			statusText := ""
			if result.Error != nil {
				statusText = errorColor(fmt.Sprintf("ERROR: %v", result.Error))
			} else if result.Status >= 200 && result.Status < 300 {
				statusText = successColor(fmt.Sprintf("OK (%d)", result.Status))
			} else {
				statusText = errorColor(fmt.Sprintf("FAIL (%d)", result.Status))
			}

			fmt.Printf("[%s] %s %s - %s - %s\n",
				methodColor(result.Method),
				serviceColor(result.ServiceName),
				// routeColor(result.RouteName),
				result.URL,
				statusText,
				timeColor(fmt.Sprintf("%dms", result.ResponseTime)),
			)

			if verbose && result.ResponseBody != "" {
				fmt.Println("Response:")
				lines := strings.Split(result.ResponseBody, "\n")
				for i, line := range lines {
					if i > 10 {
						fmt.Println("... (truncated)")
						break
					}
					fmt.Printf("  %s\n", line)
				}
				fmt.Println()
			}
		}
	}()

	// Check each route
	for _, route := range routes {
		for _, path := range route.Paths {
			methodsToCheck := []string{"GET"}

			// Use a fully qualified URL
			var baseURL string
			if len(route.HostNames) > 0 {
				baseURL = fmt.Sprintf("%s://%s", conf.DefaultProtocol, route.HostNames[0])
			} else {
				baseURL = fmt.Sprintf("%s://%s", conf.DefaultProtocol, conf.DefaultHostname)
			}

			// Build the URL
			url := baseURL
			endpoint := path
			if customPath != "" {
				if strings.HasPrefix(customPath, "/") {
					endpoint = path + customPath
				} else {
					endpoint = path + "/" + customPath
				}
			}
			url += endpoint

			for _, method := range methodsToCheck {
				wg.Add(1)
				sem <- true // Acquire token

				go func(route models.Route, url, method string) {
					defer wg.Done()
					defer func() { <-sem }() // Release token

					var result models.CheckResult
					result.ServiceName = route.ServiceName
					result.RouteName = route.Name
					result.URL = url
					result.Method = method

					start := time.Now()

					// Create request
					req, err := http.NewRequest(method, url, nil)
					if err != nil {
						result.Error = err
						result.ResponseTime = time.Since(start).Milliseconds()
						results <- result
						return
					}

					// Add headers if specified
					for key, values := range route.Headers {
						for _, value := range values {
							req.Header.Add(key, value)
						}
					}

					// Make request
					resp, err := client.Do(req)
					result.ResponseTime = time.Since(start).Milliseconds()

					if err != nil {
						result.Error = err
						results <- result
						return
					}
					defer resp.Body.Close()

					result.Status = resp.StatusCode

					// Read response body if verbose
					if verbose {
						body, err := io.ReadAll(resp.Body)
						if err != nil {
							// Just log the error but continue
							log.Printf("Error reading response body: %v", err)
						} else {
							result.ResponseBody = string(body)
						}
					}

					results <- result
				}(route, url, method)
			}
		}
	}

	wg.Wait()
	close(sem)
	close(results)

	fmt.Println("\nCheck completed.")
}
