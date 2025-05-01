package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/sahapranta/kong-check/db"
)

type ListService struct {
	*db.App
}

func NewListService(app *db.App) *ListService {
	return &ListService{
		App: app,
	}
}

func (app *ListService) ListRoutes(showHeaders bool, showMethods bool) {
	routes, err := app.GetAllRoutes()
	if err != nil {
		log.Fatalf("Failed to query routes: %v", err)
	}

	if len(routes) == 0 {
		fmt.Println("No routes found in the Kong database.")
		return
	}

	// Define colors
	serviceColor := color.New(color.FgCyan, color.Bold).SprintFunc()
	routeColor := color.New(color.FgGreen).SprintFunc()
	pathColor := color.New(color.FgYellow).SprintFunc()
	methodColor := color.New(color.FgMagenta).SprintFunc()
	headerKeyColor := color.New(color.FgBlue).SprintFunc()
	headerValColor := color.New(color.FgWhite).SprintFunc()

	// Group routes by service
	serviceRoutes := make(map[string][]string)
	for _, route := range routes {
		paths := route.Paths

		routeInfo := routeColor(fmt.Sprintf("Route: %s", route.Name))
		pathsInfo := ""

		for _, path := range paths {
			pathsInfo += fmt.Sprintf("\n    Path: %s", pathColor(path))
		}

		if showMethods && len(route.Methods) > 0 {
			methodsStr := strings.Join(route.Methods, ", ")
			pathsInfo += fmt.Sprintf("\n    Methods: %s", methodColor(methodsStr))
		}

		if showHeaders && len(route.Headers) > 0 {
			pathsInfo += "\n    Headers:"
			for key, values := range route.Headers {
				valuesStr := strings.Join(values, ", ")
				pathsInfo += fmt.Sprintf("\n      %s: %s",
					headerKeyColor(key),
					headerValColor(valuesStr))
			}
		}

		if _, exists := serviceRoutes[route.ServiceName]; !exists {
			serviceRoutes[route.ServiceName] = []string{}
		}

		serviceRoutes[route.ServiceName] = append(
			serviceRoutes[route.ServiceName],
			routeInfo+pathsInfo,
		)
	}

	// Print services and their routes
	fmt.Println("Kong API Routes: ")
	for service, routes := range serviceRoutes {
		fmt.Printf("%s\n", serviceColor(fmt.Sprintf("Service: %s", service)))
		for _, routeInfo := range routes {
			fmt.Printf("  %s\n", routeInfo)
		}
		fmt.Println()
	}

	fmt.Printf("Total: %d routes across %d services\n",
		len(routes),
		len(serviceRoutes),
	)
}
