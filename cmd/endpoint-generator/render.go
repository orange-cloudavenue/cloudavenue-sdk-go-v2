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
	"embed"
	"os"
	"strings"
	"text/template"
)

//go:embed generator.tmpl
var tmplFile embed.FS

func renderEndpoints(fileASTPackage string, endpts endpoints, config endpointGeneratorConfig) error {
	endpointTmpl := struct {
		PackageName string
		Endpoints   endpoints
	}{
		PackageName: fileASTPackage,
		Endpoints:   endpts,
	}

	tmpl, err := template.ParseFS(tmplFile, "generator.tmpl")
	if err != nil {
		return err
	}

	outputFile, err := os.Create(resolveOutputPath(config))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return tmpl.Execute(outputFile, endpointTmpl)
}

func resolveOutputPath(config endpointGeneratorConfig) string {
	pwd, _ := os.Getwd()
	split := strings.Split(pwd+"/"+config.path, "/")
	outputDir := ""

	for i := len(split) - 1; i >= 0; i-- {
		if split[i] == "cloudavenue-sdk-go-v2" {
			outputDir = strings.Join(split[:i+1], "/") + "/endpoints/"
			break
		}
	}

	if config.filename != "" {
		return outputDir + config.filename
	}

	return outputDir + "zz_" + config.output + ".go"
}
