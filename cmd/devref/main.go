/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package main

import (
	"log"
	"reflect"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"

	// Force import of all commands to register them
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
)

// NewRegistry get the global command registry
var reg = commands.NewRegistry()

func main() {

	// * Commands
	var funcs = make(map[string]Functionality)

	for _, ns := range reg.GetNamespaces() {
		funcs[ns] = loopNamespace(ns)
	}

	log.Default().Println("Found", len(funcs), "functionalities")

	// Output the functionalities to a JSON file
	err := writeJSONFile("functionalities.json", funcs)
	if err != nil {
		log.Default().Println("Error writing functionalities to JSON file:", err)
		return
	}

	// * Consoles

	var cs = make(map[string]consoles.Console)

	for _, console := range consoles.GetConsoles() {
		cs[console.SiteName] = console
	}

	// Output the functionalities to a JSON file
	err = writeJSONFile("consoles.json", cs)
	if err != nil {
		log.Default().Println("Error writing functionalities to JSON file:", err)
		return
	}

}

func loopNamespace(ns string) Functionality {
	funct := Functionality{
		Title:            ns,
		Commands:         make(map[string]Func),
		SubFunctionality: make(map[string]Functionality),
	}

	nsCmd := reg.GetCommandsByFilter(func(cmd commands.Command) bool {
		return cmd.GetNamespace() == ns && cmd.GetResource() == "" && cmd.GetVerb() == ""
	})

	if len(nsCmd) == 0 {
		log.Default().Println("No main commands reference found for namespace:", ns)
		return funct
	}

	funct.MarkdownDocumentation = nsCmd[0].MarkdownDocumentation

	// Get all commands for the namespace
	commandsByNamespace := reg.GetCommandsByFilter(func(cmd commands.Command) bool {
		return cmd.GetNamespace() == ns && cmd.GetResource() == "" && cmd.GetVerb() != ""
	})

	for _, cmd := range commandsByNamespace {
		funct.Commands[cmd.GetVerb()] = commandToFunc(cmd)
	}

	// Get all sub-commands for the namespace
	subCommands := reg.GetCommandsByFilter(func(cmd commands.Command) bool {
		return cmd.GetNamespace() == ns && cmd.GetResource() != "" && cmd.GetVerb() == ""
	})

	for _, cmd := range subCommands {
		log.Default().Println("Adding sub-command:", cmd.GetNamespace(), cmd.GetResource(), cmd.GetVerb())
		sc := loopSubCommand(cmd)
		funct.SubFunctionality[sc.Title] = sc
	}

	return funct
}

func loopSubCommand(cmd commands.Command) Functionality {
	funct := Functionality{
		Title:            cmd.GetResource(),
		Commands:         make(map[string]Func),
		SubFunctionality: make(map[string]Functionality),
	}

	funct.MarkdownDocumentation = cmd.MarkdownDocumentation

	// Get all commands for the sub-command
	commandsBySubCommand := reg.GetCommandsByFilter(func(c commands.Command) bool {
		return c.GetNamespace() == cmd.GetNamespace() && c.GetResource() == cmd.GetResource() && c.GetVerb() != ""
	})

	for _, c := range commandsBySubCommand {
		funct.Commands[c.GetVerb()] = commandToFunc(c)
	}

	return funct
}

func commandToFunc(cmd commands.Command) Func {
	f := Func{
		Namespace: cmd.GetNamespace(),
		Resource:  cmd.GetResource(),
		Verb:      cmd.GetVerb(),

		ShortDocumentation:    cmd.ShortDocumentation,
		LongDocumentation:     cmd.LongDocumentation,
		MarkdownDocumentation: cmd.MarkdownDocumentation,
	}

	// * Param
	if cmd.ParamsType != nil {
		f.Params = make([]FuncParam, 0)

		for _, spec := range cmd.ParamsSpecs {
			fType, err := commands.GetParamType(reflect.TypeOf(cmd.ParamsType), spec.Name)
			if err != nil {
				log.Default().Println("Error getting param type for", cmd.GetNamespace(), cmd.GetResource(), cmd.GetVerb(), ":", err)
				continue
			}

			fValidatorsDescription := ""
			if spec.Validators != nil {
				// If the spec has validators, we can use them to generate the description
				// e.g. "Must be a valid email address"
				for i, v := range spec.Validators {
					if v.GetMarkdownDescription() == "" {
						continue
					}
					fValidatorsDescription += v.GetMarkdownDescription()
					if i != len(spec.Validators)-1 {
						fValidatorsDescription += ", "
					} else {
						fValidatorsDescription += ". \n"
					}
				}
			}

			f.Params = append(f.Params, FuncParam{
				Name:                  spec.Name,
				Description:           spec.Description,
				Type:                  fType.String(),
				Required:              spec.Required,
				Example:               spec.Example,
				ValidatorsDescription: fValidatorsDescription,
			})
		}
	}

	// * Model
	if cmd.ModelType != nil {
		fType := reflect.TypeOf(cmd.ModelType)
		if fType.Kind() == reflect.Ptr {
			fType = fType.Elem()
		}

		docs, err := commands.GetModelTypes(fType)
		if err != nil {
			log.Default().Println("Error getting model type for", cmd.GetNamespace(), cmd.GetResource(), cmd.GetVerb(), ":", err)
			return f
		}

		f.Model = docs
	}

	// * Rules
	if cmd.ParamsRules != nil {
		exported := make([]RuleExport, 0, len(cmd.ParamsRules))
		for _, r := range cmd.ParamsRules {
			exported = append(exported, RuleExport{
				Consoles:    getConsoleNames(r.Consoles),
				WhenHuman:   ConditionToString(r.When),
				When:        ExportCondition(r.When),
				Target:      r.Target,
				Min:         r.Rule.Min,
				Max:         r.Rule.Max,
				Enum:        r.Rule.Enum,
				Pattern:     r.Rule.Pattern,
				Description: r.Rule.Description,
				Unit:        r.Rule.Unit,
			})
		}

		if len(exported) > 0 {
			f.Rules = exported
		}
	}

	return f
}
