package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sahapranta/kong-check/cmd"
	"github.com/sahapranta/kong-check/config"
	"github.com/sahapranta/kong-check/db"
	"github.com/sahapranta/kong-check/utils"
)

func main() {
	conf := config.NewConfig()
	app, err := db.NewApp(conf)

	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)

	// List command flags
	listShowHeaders := listCmd.Bool("headers", false, "Show request headers for each route")
	listShowMethods := listCmd.Bool("methods", true, "Show HTTP methods for each route")

	// Check command flags
	checkPath := checkCmd.String("path", "", "Specific path to check (appended to route)")
	checkVerbose := checkCmd.Bool("v", false, "Verbose output")
	checkTimeout := checkCmd.Int("timeout", 5, "Timeout in seconds")
	checkAll := checkCmd.Bool("all", false, "Check all routes")
	checkServices := checkCmd.String("services", "", "Comma-separated list of service names to check")
	checkMethods := checkCmd.String("methods", "", "Comma-separated list of HTTP methods to check (e.g. GET,POST,PUT)")

	if len(os.Args) < 2 {
		utils.PrintUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		cmd.NewListService(app).ListRoutes(*listShowHeaders, *listShowMethods)
	case "check":
		checkCmd.Parse(os.Args[2:])
		cmd.NewCheckService(app).CheckRoutes(*checkPath, *checkVerbose, *checkTimeout, *checkAll, *checkServices, *checkMethods)
	case "help":
		utils.PrintUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		utils.PrintUsage()
		os.Exit(1)
	}
}
