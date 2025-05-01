package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func PrintUsage() {
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
