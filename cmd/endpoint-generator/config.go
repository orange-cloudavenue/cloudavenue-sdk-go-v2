/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package main

import "flag"

func parseEndpointGeneratorConfig() (endpointGeneratorConfig, bool) {
	flagPath := flag.String("path", "", "The path to the file to generate endpoints from")
	flagFilename := flag.String("filename", "", "The name of the file to generate")
	flagDebug := flag.Bool("debug", false, "Enable debug mode")
	flagOutput := flag.String("output", "", "The output directory for the generated files")
	flag.Parse()

	config := endpointGeneratorConfig{
		path:     *flagPath,
		filename: *flagFilename,
		output:   *flagOutput,
		debug:    *flagDebug,
	}

	return config, config.path != "" && config.output != ""
}
