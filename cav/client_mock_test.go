// SPDX-FileCopyrightText: Copyright (c) 2025 Orange
// SPDX-License-Identifier: Mozilla Public License 2.0
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.
//

package cav

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/orange-cloudavenue/common-go/generator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

const (
	mockOrg = "cav01ev01ocb0001234"
)

func newMockClient() (Client, error) {
	endpoints := GetEndpointsUncategorized()
	mux := chi.NewRouter()

	for _, ep := range endpoints {
		mux.MethodFunc(ep.Method.String(), ep.PathTemplate, defaultMockHandler(ep))
	}

	hts := httptest.NewServer(mux)
	xlogger.Debug("Mock server started", slog.String("url", hts.URL))

	nC, err := NewClient(
		mockOrg,
		WithCustomEndpoints(consoles.Services{
			IHM: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/ihm",
			},
			APIVCD: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL,
			},
			APICerberus: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL,
			},
			S3: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/s3",
			},
			Netbackup: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/netbackup",
			},
		}),
		WithCloudAvenueCredential("mockuser", "mockpassword"),
	)
	if err != nil {
		return nil, err
	}

	return nC, nil
}

// defaultMockHandler returns an http.HandlerFunc that auto-generates mock responses.
func defaultMockHandler(ep *Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if ep.Name == "SessionVmware" {
			w.Header().Add("X-VMWARE-VCLOUD-ACCESS-TOKEN", "mock-access-token")
			w.Header().Set("Content-Type", "application/json")
			resp := map[string]any{
				"org":   map[string]string{"id": "{urn:org}", "name": mockOrg},
				"site":  map[string]string{"id": "{urn:site}", "name": mockOrg},
				"roles": []string{"Organization Administrator"},
			}
			respJ, _ := json.Marshal(resp)
			w.Write(respJ)
			return
		}

		// VMware job endpoints return 202 with Location header
		if ep.Name == "GetJobVmware" {
			w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Cerberus job endpoints return 201 with jobId
		if ep.Name == "GetJobCerberus" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"jobId":"87ab1934-0146-4fb0-80bc-815fea03214d","message":"Job created successfully"}`))
			return
		}

		if ep.ResponseType != nil {
			bodyType := reflect.TypeOf(ep.ResponseType)
			if bodyType.Kind() == reflect.Pointer {
				bodyType = bodyType.Elem()
			}

			newBodyType := reflect.PointerTo(bodyType)
			newBody := reflect.New(newBodyType).Interface()

			switch bodyType.Kind() {
			case reflect.Slice:
				generator.Slice(newBody)
			default:
				if err := generator.Struct(newBody); err != nil {
					xlogger.Error("Error generating mock data", slog.String("endpoint", ep.Name), slog.Any("error", err))
				}
			}

			w.Header().Set("Content-Type", "application/json")
			respJ, err := json.Marshal(newBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(respJ)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Mock response"}`))
	}
}
