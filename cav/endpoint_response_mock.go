/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-faker/faker/v4"
)

var defaultMockResponseFunc = func(ep *Endpoint) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {

		returnErrFromStatusCodeExpected(w, ep.mockResponseStatusCode)

		var body any

		if ep.mockResponseData != nil {
			log.Default().Printf("Using mock response data for endpoint %s", ep.Name)
			// If mock response data is defined, use it directly
			body = ep.mockResponseData
		} else if ep.BodyType != nil {
			log.Default().Printf("Generating mock response data for endpoint %s", ep.Name)
			body = ep.BodyType
			if err := faker.FakeData(&body); err != nil {
				log.Default().Println("Error generating mock data for endpoint:", ep.Name, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		bodyEncoded, err := json.Marshal(body)
		if err != nil {
			log.Default().Println("Error encoding body for endpoint:", ep.Name, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Default().Printf("Mock response for endpoint %s: %s", ep.PathTemplate, string(bodyEncoded))

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyEncoded)
	}
}

var (
	GetDefaultMockResponseFunc  = defaultMockResponseFunc
	PostDefaultMockResponseFunc = defaultMockResponseFunc
)

func returnErrFromStatusCodeExpected(w http.ResponseWriter, statusCode *int) {
	if statusCode == nil {
		log.Default().Println("No status code defined for mock response, returning 200 OK")
		return
	}
	log.Default().Println("Checking status code for mock response:", *statusCode)

	if *statusCode >= 200 && *statusCode < 300 {
		return
	}

	log.Default().Printf("Mock response error for status code %d", *statusCode)
	http.Error(w, http.StatusText(*statusCode), *statusCode)
	w.WriteHeader(*statusCode)
}
