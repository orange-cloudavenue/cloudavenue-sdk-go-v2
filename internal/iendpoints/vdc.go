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
	"slices"
	"strings"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/urn"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate endpoint-generator -path vdc.go -output vdc

func init() {
	// ListVDC
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/38.1/doc/types/ReferenceType.html",
		Name:             "ListVdc",
		Description:      "List VDCs",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/api/query",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "Filter to apply to the list of VDCs. Format: key==value. Allowed keys: name, id.",
				ValidatorFunc: func(value string) error {
					valueSplit := strings.Split(value, "==")
					if len(valueSplit) != 2 {
						return errors.New("filter must be in the format 'key==value'")
					}

					allowedKeys := []string{"name", "id"}
					if !slices.Contains(allowedKeys, valueSplit[0]) {
						return fmt.Errorf("filter key '%s' is not allowed", valueSplit[0])
					}

					return nil
				},
				TransformFunc: func(value string) (string, error) {
					// Add ( ) around the filter value
					return fmt.Sprintf("(%s)", value), nil
				},
			},
			{
				Name:        "pageSize",
				Description: "The number of items per page.",
				Value:       "100",
			},
			{
				Name:        "format",
				Description: "The format of the response.",
				Value:       "records",
			},
			{
				Name:        "type",
				Description: "The type of object to query",
				Value:       "orgVdc",
			},
		},
		RequestMiddlewares: []resty.RequestMiddleware{
			func(_ *resty.Client, req *resty.Request) error {
				// Set the Accept header to application/*+json;version=38.1
				req.SetHeader("Accept", "application/*+json;version=38.1")
				return nil
			},
		},
		ResponseMiddlewares: []resty.ResponseMiddleware{
			func(_ *resty.Client, resp *resty.Response) error {
				r := resp.Result().(*itypes.ApiResponseListVDC)

				// Extract ID from HREF
				for i, record := range r.Records {
					id, err := extractor.ExtractUUID(record.HREF)
					if err != nil {
						return fmt.Errorf("failed to extract ID from HREF: %w", err)
					}
					r.Records[i].ID = urn.Normalize(urn.VDC, id).String()
				}

				return nil
			},
		},
		BodyResponseType: itypes.ApiResponseListVDC{},
		MockResponseFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := itypes.ApiResponseListVDC{
				Records: make([]itypes.ApiResponseListVDCRecord, 0),
			}

			// If QueryParam "filter" is set, return a filtered response
			if r.URL.Query().Get("filter") != "" {
				filter := r.URL.Query().Get("filter")

				filterParts := strings.Split(filter, "==")

				r := &itypes.ApiResponseListVDCRecord{}
				generator.MustStruct(r)

				r.ID = func() string {
					if filterParts[0] == "id" {
						return filterParts[1]
					}
					return ""
				}()
				r.Name = func() string {
					if filterParts[0] == "name" {
						return fmt.Sprintf("mockvdc-%s", filterParts[1])
					}
					return generator.MustGenerate("mockvdc-{word}")
				}()
				resp.Records = append(resp.Records, *r)
			} else {
				generator.MustStruct(&resp)
			}

			// json encode
			w.Header().Set("Content-Type", "application/json")
			respJ, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = w.Write(respJ)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}),
	}.Register()

	// GetVDC
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Vdc.html",
		Name:             "GetVdc",
		Description:      "Get VDC",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/api/vdc/{vdc-id}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-id",
				Description: "The ID of the VDC.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdc")
				},
				TransformFunc: func(value string) (string, error) {
					// vdc-id require UUID format and not urn format
					return extractor.ExtractUUID(value)
				},
			},
		},
		BodyResponseType: itypes.ApiResponseGetVDC{},
		RequestMiddlewares: []resty.RequestMiddleware{
			func(_ *resty.Client, req *resty.Request) error {
				// Set the Accept header to application/*+json;version=38.1
				req.SetHeader("Accept", "application/*+json;version=38.1")
				return nil
			},
		},
		MockResponseFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := &itypes.ApiResponseGetVDC{}

			generator.MustStruct(resp)

			// Extract the VDC ID from the path
			vdcID := strings.Split(r.URL.Path, "/")[5]

			// Overwrite the ID in the response
			resp.ID = urn.Normalize(urn.VDC, vdcID).String()

			// json encode
			w.Header().Set("Content-Type", "application/json")
			respJ, _ := json.Marshal(resp)
			_, _ = w.Write(respJ)
		}),
	}.Register()

	// GetVDCMetadata
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VdcMetadata.html",
		Name:             "GetVdcMetadata",
		Description:      "Get VDC Metadata",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/api/vdc/{vdc-id}/metadata",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-id",
				Description: "The ID of the VDC.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdc")
				},
				TransformFunc: func(value string) (string, error) {
					// vdc-id require UUID format and not urn format
					return extractor.ExtractUUID(value)
				},
			},
		},
		BodyResponseType: itypes.ApiResponseGetVDCMetadatas{},
		RequestMiddlewares: []resty.RequestMiddleware{
			func(_ *resty.Client, req *resty.Request) error {
				// Set the Accept header to application/*+json;version=38.1
				req.SetHeader("Accept", "application/*+json;version=38.1")
				return nil
			},
		},
		MockResponseFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := &itypes.ApiResponseGetVDCMetadatas{
				Metadatas: make([]itypes.ApiResponseGetVDCMetadata, 0),
			}

			resp.Metadatas = append(resp.Metadatas,
				itypes.ApiResponseGetVDCMetadata{
					Name:  "vdcBillingModel",
					Value: itypes.ApiResponseGetVDCMetadataValue{Value: "PAYG"},
				},
				itypes.ApiResponseGetVDCMetadata{
					Name:  "vdcStorageBillingModel",
					Value: itypes.ApiResponseGetVDCMetadataValue{Value: "PAYG"},
				},
				itypes.ApiResponseGetVDCMetadata{
					Name:  "vdcServiceClass",
					Value: itypes.ApiResponseGetVDCMetadataValue{Value: "HP"},
				},
				itypes.ApiResponseGetVDCMetadata{
					Name:  "vdcDisponibilityClass",
					Value: itypes.ApiResponseGetVDCMetadataValue{Value: "ONE-ROOM"},
				})

			// json encode
			w.Header().Set("Content-Type", "application/json")
			respJ, _ := json.Marshal(resp)

			_, _ = w.Write(respJ)
		}),
	}.Register()

	// CreateVdc
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/createOrgVdc",
		Name:             "CreateVdc",
		Description:      "Create a new Org VDC",
		Method:           cav.MethodPOST,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/vdcs",
		BodyRequestType:  itypes.ApiRequestCreateVDC{},
		BodyResponseType: cav.Job{},
	}.Register()

	// UpdateVdc
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/updateOrgVdc",
		Name:             "UpdateVdc",
		Description:      "Update an existing Org VDC",
		Method:           cav.MethodPUT,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/vdcs/{vdc-name}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-name",
				Description: "The name of the VDC to update.",
				Required:    true,
			},
		},
		BodyRequestType:  itypes.ApiRequestUpdateVDC{},
		BodyResponseType: cav.Job{},
	}.Register()

	// DeleteVdc
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/deleteOrgVdc",
		Name:             "DeleteVdc",
		Description:      "Delete an existing Org VDC",
		Method:           cav.MethodDELETE,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/vdcs/{vdc-name}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-name",
				Description: "The name of the VDC to delete.",
				Required:    true,
			},
		},
		BodyResponseType: cav.Job{},
	}.Register()
}
