/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/urn"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

//go:generate endpoint-generator -path storage_profile.go -output storage_profile

func init() {
	// * ListStorageProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/queries/orgVdcStorageProfile.html",
		Name:             "ListStorageProfile",
		Description:      "List VDC Storage Profiles",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/api/query/",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "ID or Name of the VDC to get storage profiles for.",
				ValidatorFunc: func(value string) error {
					// Get the filter value
					valueSplit := strings.Split(value, "==")
					if len(valueSplit) != 2 {
						return errors.New("filter must be in the format 'key==value'")
					}

					// Validate the value based on the key
					switch valueSplit[0] {
					case "vdc": //nolint: goconst
						// vdc require an urn format for VDC
						return validators.New().Var(value, "urn=vdc")
					case "vdcName", "name": //nolint: goconst
						// No specific format required
						return nil
					case "id":
						// ID require an urn format for VDC Storage Profile
						return validators.New().Var(value, "urn=vdcstorageProfile")
					default:
						// If the key is not recognized, return an error
						return fmt.Errorf("filter key '%s' is not allowed", valueSplit[0])
					}
				},
				TransformFunc: func(value string) (string, error) {
					// Get the filter value
					valueSplit := strings.Split(value, "==")

					// Extract the value based on the key
					switch valueSplit[0] {
					case "vdc":
						// vdc-id require UUID format and not urn format
						v, err := extractor.ExtractUUID(valueSplit[1])
						if err != nil {
							return "", err
						}
						return fmt.Sprintf("vdc==%s", v), nil
					case "vdcName":
						// vdcName is a string, no specific format required
						v := valueSplit[1]
						return fmt.Sprintf("vdcName==%s", v), nil
					case "name":
						// Name is a string, no specific format required
						v := valueSplit[1]
						return fmt.Sprintf("name==%s", v), nil
					case "id":
						// ID require an urn format for VDC Storage Profile
						v, err := extractor.ExtractUUID(valueSplit[1])
						if err != nil {
							return "", err
						}
						return fmt.Sprintf("id==%s", v), nil
					default:
						// If the key is not recognized, return an error
						return "", fmt.Errorf("filter key '%s' is not allowed", valueSplit[0])
					}
				},
			},
			{
				Name:        "pageSize",
				Description: "The number of items per page.",
				Value:       "30",
			},
			{
				Name:        "format",
				Description: "The format of the response.",
				Value:       "records",
			},
			{
				Name:        "type",
				Description: "The type of object to query",
				Value:       "orgVdcStorageProfile",
			},
			{
				Name:        "sortAsc",
				Description: "Sort the results in ascending order.",
				Value:       "name",
			},
		},
		RequestMiddlewares: []resty.RequestMiddleware{
			func(_ *resty.Client, req *resty.Request) error {
				// Set the Accept header to application/*+json;version=38.1
				req.SetHeader("Accept", "application/*+json;version=38.1")
				return nil
			},
		},
		BodyResponseType: apiResponseListStorageProfiles{},
		// ResponseMiddlewares is used to extract the ID from the response and set it in the context
		ResponseMiddlewares: []resty.ResponseMiddleware{
			func(_ *resty.Client, resp *resty.Response) error {
				r := resp.Result().(*apiResponseListStorageProfiles)

				for i, strPro := range r.StorageProfiles {
					// Extract ID from HREF
					id, err := extractor.ExtractUUID(strPro.HREF)
					if err != nil {
						return fmt.Errorf("failed to extract ID from HREF: %w", err)
					}
					r.StorageProfiles[i].ID = urn.Normalize(urn.VDCStorageProfile, id).String()

					// Extract VDC ID from HREF
					vdcID, err := extractor.ExtractUUID(strPro.VdcId)
					if err != nil {
						return fmt.Errorf("failed to extract VDC ID from HREF: %w", err)
					}
					r.StorageProfiles[i].VdcId = urn.Normalize(urn.VDC, vdcID).String()
				}

				return nil
			},
		},

		MockResponseFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := &apiResponseListStorageProfiles{
				StorageProfiles: make([]apiResponseListStorageProfile, 0),
			}

			// If QueryParam "filter" is set, return a filtered response
			if r.URL.Query().Get("filter") != "" {
				filter := r.URL.Query().Get("filter")
				filterParts := strings.Split(filter, "==")

				r := &apiResponseListStorageProfile{}
				generator.MustStruct(r)

				r.Name = func() string {
					if filterParts[0] == "name" {
						return filterParts[1]
					}
					return "platinum3k_r1"
				}()

				r.ID = func() string {
					if filterParts[0] == "id" {
						return filterParts[1]
					}
					return ""
				}()

				r.VdcId = func() string {
					if filterParts[0] == "vdc" {
						return urn.Normalize(urn.VDC, filterParts[1]).String()
					}
					return generator.MustGenerate("{urn:vdc}")
				}()

				r.VdcName = func() string {
					if filterParts[0] == "vdcName" {
						return filterParts[1]
					}
					return generator.MustGenerate("{word}")
				}()
				resp.StorageProfiles = append(resp.StorageProfiles, *r)
			} else {
				// If no filter is set, generate a random response
				generator.MustStruct(resp)
			}

			// json encode
			w.Header().Set("Content-Type", "application/json")
			respJ, _ := json.Marshal(resp)

			_, _ = w.Write(respJ)
		}),
	}.Register()
}
