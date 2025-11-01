package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"

	bboxclient "bbox-cli/client"

	"github.com/joho/godotenv"
)

func Run() {
	godotenv.Load()
	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(1)
	}

	// Get password from env
	password := os.Getenv("BBOX_PWD")
	if password == "" {
		fmt.Println("Error: BBOX_PWD env variable not set")
		os.Exit(1)
	}

	// Parse URL
	baseURL := "https://mabbox.bytel.fr/api/v1"
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	// Create client
	client, err := bboxclient.NewClient(parsedURL)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Authenticate
	authInterface := client.Auth()
	if err := authInterface.BasicAuth(password); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Start token refresher
	authInterface.StartTokenRefresher()

	// Parse subcommand
	subcommand := os.Args[1]

	switch subcommand {
	case "nat":
		handleNat(client, os.Args[2:])
	case "firewall":
		handleFirewall(client, os.Args[2:])
	case "help":
		PrintUsage()
	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		PrintUsage()
		os.Exit(1)
	}
}

func PrintUsage() {
	fmt.Println("bboxcli - Bbox Configuration Tool")
	fmt.Println()
	fmt.Println("Usage: bboxcli <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  firewall show        Show all firewall rules")
	fmt.Println("  firewall show <id>   Show detailed firewall rule")
	fmt.Println("  firewall add <rule>  Add a new firewall rule")
	fmt.Println("  firewall delete <id> Delete a firewall rule")
	fmt.Println("  nat show             Show all NAT rules")
	fmt.Println("  nat show <id>        Show detailed NAT rule")
	fmt.Println("  nat enable <id>      Enable a NAT rule")
	fmt.Println("  nat disable <id>     Disable a NAT rule")
	fmt.Println("  help                 Show this help message")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  BBOX_PWD            Password for Bbox authentication (required, can be set in .env file)")
}
