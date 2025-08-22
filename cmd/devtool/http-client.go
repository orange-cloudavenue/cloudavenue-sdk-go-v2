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
	"context"
	"fmt"
	"strings"

	"github.com/niemeyer/pretty"
	"github.com/spf13/cobra"
	"resty.dev/v3"
)

type clientRawRequest interface {
	NewRawRequest(ctx context.Context, subclientName string) (req *resty.Request, err error)
}

var endpointHTTP = &cobra.Command{
	Use:   "http [name]",
	Short: "Test an http Endpoint with custom subclient",
	Long:  `Allows you to test an HTTP endpoint by specifying its custom subclient.`,
	Args:  cobra.ExactArgs(1),
	Example: `
# GET
http --subclient vmware --logger debug --header "Accept:application/*;version=38.1" /cloudapi/1.0.0/site/settings

# POST
http --subclient vmware --logger debug --body='{"key":"value"}' --method=POST /my/endpoint
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := newClient()
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			return
		}

		clientRaw, ok := client.(clientRawRequest)
		if !ok {
			fmt.Printf("Error asserting client to clientRawRequest: %v\n", err)
			return
		}

		req, err := clientRaw.NewRawRequest(context.Background(), cmdHttpSubclientName)
		if err != nil {
			fmt.Printf("Error creating raw request: %v\n", err)
			return
		}

		req.SetBody(cmdHttpBody)
		req.EnableDebug()
		req.EnableTrace()

		for _, header := range cmdHttpHeaders {
			h := strings.SplitN(header, ":", 2)
			if len(h) != 2 {
				fmt.Printf("Invalid header format: %s\n", header)
				continue
			}
			req.SetHeader(strings.TrimSpace(h[0]), strings.TrimSpace(h[1]))
		}

		resp, err := req.Execute(cmdHttpMethod, args[0])
		if err != nil {
			fmt.Printf("Error executing request: %v\n", err)
			return
		}

		pretty.Print(resp.Result())
	},
}

var (
	cmdHttpSubclientName string
	cmdHttpBody          string
	cmdHttpMethod        string
	cmdHttpHeaders       []string
)

func init() {
	endpointHTTP.Flags().StringVar(&cmdHttpBody, "body", "", "Request body as a JSON string")
	endpointHTTP.Flags().StringVar(&cmdHttpSubclientName, "subclient", "", "Subclient name")
	endpointHTTP.Flags().StringVar(&cmdHttpMethod, "method", "GET", "HTTP method")
	endpointHTTP.Flags().StringSliceVar(&cmdHttpHeaders, "header", []string{}, "Request headers (key:value)")

	rootCmd.AddCommand(endpointHTTP)
}
