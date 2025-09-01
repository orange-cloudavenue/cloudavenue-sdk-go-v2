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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/urn"
	"github.com/orange-cloudavenue/common-go/validators"
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
				Description: "Filter to apply to the list of VDC Storage Profile. Format: key==value. Supported keys: vdc, vdcName, name, id.",
				ValidatorFunc: func(value string) error {
					// Support multiple filters separated by ';'
					filters := strings.Split(value, ";")
					for _, filter := range filters {
						valueSplit := strings.Split(filter, "==")
						if len(valueSplit) != 2 {
							return errors.New("filter must be in the format 'key==value' or 'key1==value1;key2==value2'")
						}
						switch valueSplit[0] {
						case "vdc":
							if err := validators.New().Var(valueSplit[1], "urn=vdc"); err != nil {
								return err
							}
						case "vdcName", "name":
							// No specific format required
						case "id":
							if err := validators.New().Var(valueSplit[1], "urn=vdcstorageProfile"); err != nil {
								return err
							}
						default:
							return fmt.Errorf("filter key '%s' is not allowed", valueSplit[0])
						}
					}
					return nil
				},
				TransformFunc: func(value string) (string, error) {
					// Support multiple filters separated by ';'
					filters := strings.Split(value, ";")
					var transformed []string
					for _, filter := range filters {
						valueSplit := strings.Split(filter, "==")
						if len(valueSplit) != 2 {
							return "", errors.New("filter must be in the format 'key==value' or 'key1==value1;key2==value2'")
						}
						switch valueSplit[0] {
						case "vdc":
							v, err := extractor.ExtractUUID(valueSplit[1])
							if err != nil {
								return "", err
							}
							transformed = append(transformed, fmt.Sprintf("vdc==%s", v))
						case "vdcName":
							v := valueSplit[1]
							transformed = append(transformed, fmt.Sprintf("vdcName==%s", v))
						case "name":
							v := valueSplit[1]
							transformed = append(transformed, fmt.Sprintf("name==%s", v))
						case "id":
							v, err := extractor.ExtractUUID(valueSplit[1])
							if err != nil {
								return "", err
							}
							transformed = append(transformed, fmt.Sprintf("id==%s", v))
						default:
							return "", fmt.Errorf("filter key '%s' is not allowed", valueSplit[0])
						}
					}
					return strings.Join(transformed, ";"), nil
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
		BodyResponseType: itypes.ApiResponseListStorageProfiles{},
		// ResponseMiddlewares is used to extract the ID from the response and set it in the context
		ResponseMiddlewares: []resty.ResponseMiddleware{
			func(_ *resty.Client, resp *resty.Response) error {
				// Extract ID and VDCID from HREF and set it in the context
				r := resp.Result().(*itypes.ApiResponseListStorageProfiles)

				for i, strPro := range r.StorageProfiles {
					// Extract ID from HREF
					id, err := extractor.ExtractUUID(strPro.HREF)
					if err != nil {
						return fmt.Errorf("failed to extract ID from HREF: %w", err)
					}
					r.StorageProfiles[i].ID = urn.Normalize(urn.VDCStorageProfile, id).String()

					// Extract VDC ID from HREF
					vdcID, err := extractor.ExtractUUID(strPro.VdcID)
					if err != nil {
						return fmt.Errorf("failed to extract VDC ID from HREF: %w", err)
					}
					r.StorageProfiles[i].VdcID = urn.Normalize(urn.VDC, vdcID).String()
				}

				return nil
			},
			// modify response for limit and used fields to transform values in GiB instead of MiB
			func(_ *resty.Client, resp *resty.Response) error {
				r := resp.Result().(*itypes.ApiResponseListStorageProfiles)

				for i, strPro := range r.StorageProfiles {
					r.StorageProfiles[i].Limit = strPro.Limit / 1024
					r.StorageProfiles[i].Used = strPro.Used / 1024
				}

				return nil
			},
		},

		MockResponseFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: make([]itypes.ApiResponseListStorageProfile, 0),
			}

			// If QueryParam "filter" is set, return a filtered response
			if r.URL.Query().Get("filter") != "" {
				filter := r.URL.Query().Get("filter")
				filterParts := strings.Split(filter, "==")

				r := &itypes.ApiResponseListStorageProfile{}
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

				r.VdcID = func() string {
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
