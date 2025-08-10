package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

var helpCmd = &cobra.Command{
	Use:   "help [endpoint|command] [name]",
	Short: "Show help for an endpoint or command",
	Long:  `Display documentation for a specific endpoint or CLI command.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		kind := strings.ToLower(args[0])
		name := args[1]

		switch kind {
		case "endpoint":
			showEndpointHelp(name)
		// case "command":
		// 	showCommandHelp(name)
		default:
			fmt.Printf("Unknown help kind: %s. Use 'endpoint' or 'command'.\n", kind)
			os.Exit(1)
		}
	},
}

func showEndpointHelp(name string) {
	var ep *cav.Endpoint
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Endpoint '%s' not found.\n", name)
			os.Exit(1)
		}
	}()
	ep = cav.MustGetEndpoint(name)

	fmt.Printf("Name: %s\n", ep.Name)
	fmt.Printf("Description: %s\n", ep.Description)
	fmt.Printf("Documentation: %s\n", ep.DocumentationURL)
	fmt.Printf("Method: %s\n", ep.Method)
	fmt.Printf("Path: %s\n", ep.PathTemplate)

	if len(ep.PathParams) > 0 {
		fmt.Println("Path Parameters:")
		for _, p := range ep.PathParams {
			req := ""
			if p.Required {
				req = " [Required]"
			}
			fmt.Printf("  - %s%s: %s\n", p.Name, req, p.Description)
		}
	}
	if len(ep.QueryParams) > 0 {
		fmt.Println("Query Parameters:")
		for _, q := range ep.QueryParams {
			req := ""
			if q.Required {
				req = " [Required]"
			}
			fmt.Printf("  - %s%s: %s\n", q.Name, req, q.Description)
		}
	}
}

// func showCommandHelp(name string) {
// 	reg := commands.NewRegistry()
// 	cmd, err := reg.SearchCommands(name)
// 	if err != nil {
// 		fmt.Printf("Command '%s' not found.\n", name)
// 		os.Exit(1)
// 	}

// 	fmt.Printf("Command: %s\n", cmd.GetNamespace())
// 	fmt.Printf("Resource: %s\n", cmd.GetResource())
// 	fmt.Printf("Verb: %s\n", cmd.GetVerb())
// 	fmt.Printf("Description: %s\n", cmd.LongDocumentation)
// }

func init() {
	rootCmd.AddCommand(helpCmd)
}
