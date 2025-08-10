package main

import (
	"context"
	"reflect"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/niemeyer/pretty"
	"github.com/spf13/cobra"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
)

var (
	commandParams map[string]string
	namespace     string
	resource      string
	verb          string
)

var commandCmd = &cobra.Command{
	Use:                "command [namespace] [<resource>] [verb]",
	Short:              "Test a Command with custom parameters",
	Long:               `Allows you to test a CloudAvenue Command by specifying its name and parameter values.`,
	Args:               cobra.MinimumNArgs(2),
	DisableFlagParsing: true, // Disable flag parsing to allow custom parameters
	Run: func(cmd *cobra.Command, args []string) {
		commandParams = make(map[string]string)

		namespace = args[0]

		if strings.HasPrefix(args[2], "--") {
			verb = args[1]
		} else {
			resource = args[1]
			verb = args[2]
		}

		for i, entry := range args {
			if strings.HasPrefix(entry, "--") {
				key := strings.TrimPrefix(entry, "--")
				commandParams[key] = args[i+1]
			}
		}

		reg := commands.NewRegistry()
		cmds := reg.GetCommandsByFilter(func(cmd commands.Command) bool {
			return strings.EqualFold(cmd.GetNamespace(), namespace) && strings.EqualFold(cmd.GetResource(), resource) && strings.EqualFold(cmd.GetVerb(), verb)
		})
		if len(cmds) == 0 {
			log.Error("Command not found", "namespace", namespace, "resource", resource, "verb", verb)
			return
		}

		command := cmds[0]

		log.Info("Executing command", "namespace", command.GetNamespace(), "resource", command.GetResource(), "verb", command.GetVerb())

		rType := reflect.TypeOf(command.ParamsType)
		rVal := reflect.New(rType).Elem()

		for paramName, paramValue := range commandParams {
			if err := commands.StoreValueAtPath(rVal.Addr().Interface(), paramName, paramValue); err != nil {
				log.Error("Error storing parameter value", "param", paramName, "error", err)
				return
			}
		}

		log.Info("Parameters set")
		pretty.Print(rVal.Interface())

		client, err := newClient()
		if err != nil {
			log.Error("Error creating client", "error", err)
			return
		}

		vdcClient, err := vdc.New(client)
		if err != nil {
			log.Error("Error creating VDC client", "error", err)
			return
		}

		// Call the command's RunnerFunc if defined
		if command.RunnerFunc == nil {
			log.Error("No runner function defined for this command")
			return
		}
		log.Info("Running command")
		result, err := command.Run(context.Background(), vdcClient, rVal.Interface())
		if err != nil {
			log.Error("Error executing command", "error", err)
			return
		}

		pretty.Print(result)
	},
}

func init() {
	rootCmd.AddCommand(commandCmd)
}
