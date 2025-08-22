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
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/niemeyer/pretty"
	"github.com/spf13/cobra"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

var (
	pathParams  []string
	queryParams []string
	bodyJSON    string
)

var endpointCmd = &cobra.Command{
	Use:   "endpoint [name]",
	Short: "Test an Endpoint with custom parameters",
	Long:  `Allows you to test a CloudAvenue Endpoint by specifying its name, path parameters, query parameters, and request body.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpointName := args[0]
		client, err := newClient()
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			return
		}
		ep := cav.MustGetEndpoint(endpointName)

		// QueryParams
		queryParamsMap := make(map[*cav.QueryParam]string)
		for _, qp := range queryParams {
			parts := strings.SplitN(qp, ":", 2)
			if len(parts) != 2 {
				logger.Error("Invalid query parameter format", "param", qp)
				return
			}
			for _, epQP := range ep.QueryParams {
				log.Debug("Checking query parameter", "name", epQP.Name, "expected", parts[0])
				if epQP.Name == parts[0] {
					queryParamsMap[&epQP] = parts[1]
				}
			}
		}

		// PathParams
		pathParamsMap := make(map[*cav.PathParam]string)
		for _, pp := range pathParams {
			parts := strings.SplitN(pp, ":", 2)
			if len(parts) != 2 {
				logger.Error("Invalid path parameter format", "param", pp)
				return
			}
			logger.Debug("Processing path parameter", "name", parts[0], "value", parts[1])
			for _, epPP := range ep.PathParams {
				if epPP.Name == parts[0] {
					pathParamsMap[&epPP] = parts[1]
				}
			}
		}

		// Body
		var body any
		if bodyJSON != "" {
			bodyType := ep.BodyRequestType
			if bodyType == nil {
				logger.Error("This endpoint does not accept a body")
				return
			}
			bodyPtr := reflect.New(reflect.TypeOf(bodyType)).Interface()
			if err := json.Unmarshal([]byte(bodyJSON), bodyPtr); err != nil {
				logger.Error("Failed to unmarshal body JSON", "error", err)
				return
			}
			body = reflect.ValueOf(bodyPtr).Elem().Interface()
		}
		cancel := spinner("Waiting...", monkeys, 200*time.Millisecond)
		resp, err := client.Do(
			context.Background(),
			ep,
			func() []cav.EndpointRequestOption {
				var opts []cav.EndpointRequestOption
				for param, value := range pathParamsMap {
					opts = append(opts, cav.WithPathParam(*param, value))
				}
				for param, value := range queryParamsMap {
					opts = append(opts, cav.WithQueryParam(*param, value))
				}
				if body != nil {
					opts = append(opts, cav.SetBody(body))
				}
				opts = append(opts,
					cav.SetCustomRestyOption(func(req *resty.Request) {
						req.SetDebug(true)
						req.SetTrace(true)
					}),
				)
				return opts
			}()...,
		)
		cancel()
		if err != nil {
			logger.Error("Error executing endpoint", "endpoint", endpointName, "error", err)
			return
		}

		pretty.Print(resp.Result())
	},
}

func init() {
	endpointCmd.Flags().StringArrayVarP(&pathParams, "path-param", "p", []string{}, "Path parameter in the form key:value (can be specified multiple times)")
	endpointCmd.Flags().StringArrayVarP(&queryParams, "query-param", "q", []string{}, "Query parameter in the form key:value (can be specified multiple times)")
	endpointCmd.Flags().StringVar(&bodyJSON, "body", "", "Request body as a JSON string")

	rootCmd.AddCommand(endpointCmd)
}
