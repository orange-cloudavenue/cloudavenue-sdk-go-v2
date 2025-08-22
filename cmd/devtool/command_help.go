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
	"html/template"
	"os"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
)

func help(cmd commands.Command) {

	tmpl := `
{{ .Cmd.LongDocumentation }}
 
Usage:
    {{ .Cmd.Namespace | ToLower }} {{ .Cmd.Resource | ToLower }} {{ .Cmd.Verb | ToLower }} {{ if .Cmd.ParamsSpecs }}[OPTIONS]{{ end }}

{{ if .Cmd.ParamsSpecs }}
Options:
{{- range $index, $element := .Cmd.ParamsSpecs }}
    --{{ $element.Name }}
        {{ $element.Description }}
{{- end -}}
{{- else if .SubCommands }}
Subcommands:
{{- range $index, $element := .SubCommands }}
	{{- if and (ne $element.Resource "") (ne $element.Verb "") }}
    {{ $element.Resource | ToLower }} {{ $element.Verb | ToLower }} -- {{ $element.ShortDocumentation }}
	{{- end -}}
	{{- if and (eq $element.Resource "") (ne $element.Verb "") }}
    {{ $element.Verb | ToLower }} -- {{ $element.ShortDocumentation }}
	{{- end -}}
{{- end -}}
{{- end -}}
`

	// Find if command have subcommands
	subCommands := commands.NewRegistry().GetCommandsByFilter(func(c commands.Command) bool {
		if cmd.GetResource() != "" {
			return cmd.GetNamespace() == c.GetNamespace() && c.GetResource() == cmd.GetResource()
		} else {
			return cmd.GetNamespace() == c.GetNamespace()
		}
	})
	if len(subCommands) == 0 {
		// No subcommands found
	}

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	templateData := struct {
		Cmd         commands.Command
		SubCommands []commands.Command
	}{
		Cmd:         cmd,
		SubCommands: subCommands,
	}

	// Use go template
	t, err := template.New("help").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, templateData)
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
