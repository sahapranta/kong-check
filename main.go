package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sahapranta/kong-check/cmd"
	"github.com/sahapranta/kong-check/config"
)

func main() {
	conf := config.NewConfig()

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

	// Handle commands
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		cmd.ListRoutes(conf, *listShowHeaders, *listShowMethods)
	case "check":
		checkCmd.Parse(os.Args[2:])
		cmd.CheckRoutes(conf, *checkPath, *checkVerbose, *checkTimeout, *checkAll, *checkServices, *checkMethods)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Println(bold("Kong API Explorer"))
	fmt.Println("\nUsage:")
	fmt.Printf("  %s [command] [options]\n\n", os.Args[0])

	fmt.Println(bold("Commands:"))
	fmt.Printf("  %s\t\tList all routes from Kong\n", green("list"))
	fmt.Printf("    %s\tShow request headers for each route\n", green("--headers"))
	fmt.Printf("    %s\tShow HTTP methods for each route (default: true)\n", green("--methods"))

	fmt.Printf("  %s\t\tCheck route endpoints using HTTP requests\n", green("check"))
	fmt.Printf("    %s PATH\tSpecific path to append to routes\n", green("--path"))
	fmt.Printf("    %s\t\tVerbose output including response details\n", green("-v"))
	fmt.Printf("    %s N\tTimeout in seconds (default: 5)\n", green("--timeout"))
	fmt.Printf("    %s\t\tCheck all routes\n", green("--all"))
	fmt.Printf("    %s NAMES\tComma-separated list of service names to check\n", green("--services"))
	fmt.Printf("    %s METHODS\tComma-separated list of HTTP methods to check (e.g. GET,POST,PUT)\n", green("--methods"))

	fmt.Printf("  %s\t\tShow this help message\n", green("help"))

	fmt.Println("\nExamples:")
	fmt.Println("  List all routes:")
	color.RGB(150, 150, 150).Printf("    %s list\n", os.Args[0])
	fmt.Println("  Check health of all routes:")
	color.RGB(150, 150, 150).Printf("    %s check --all\n", os.Args[0])
	fmt.Println("  Check specific services with a custom path:")
	color.RGB(150, 150, 150).Printf("    %s check --services auth-service,user-service --path /health\n", os.Args[0])
	fmt.Println("  Check routes with specific HTTP methods:")
	color.RGB(150, 150, 150).Printf("    %s check --all --methods GET,POST\n", os.Args[0])
}
