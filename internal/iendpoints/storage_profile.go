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
	"errors"
	"fmt"
	"strings"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path storage_profile.go -output storage_profile

func init() {
	// * ListStorageProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/queries/orgVdcStorageProfile.html",
		Name:             "ListStorageProfile",
		Description:      "List VDC Storage Profiles",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/query/",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "Filter to apply to the list of VDC Storage Profile. Format: key==value. Supported keys: vdc, vdcName, name, id.",
				ValidatorFunc: func(value string) error {
					// Support multiple filters separated by ';'
					filters := strings.SplitSeq(value, ";")
					for filter := range filters {
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
		ResponseType: itypes.ApiResponseListStorageProfiles{},
	}.Register()
}
