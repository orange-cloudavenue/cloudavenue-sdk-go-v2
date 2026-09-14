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

func main() {
	config, ok := parseEndpointGeneratorConfig()
	if !ok {
		flag.Usage()
		return
	}

	fileAST, endpts, err := scanEndpoints(config.path, config.debug)
	if err != nil {
		panic(err)
	}

	err = renderEndpoints(fileAST.Name.Name, endpts, config)
	if err != nil {
		panic(err)
	}
}
