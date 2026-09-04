/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package inetworkobjects

import (
	"context"
	"fmt"
	"regexp"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

var appPortProfileScopes = []string{
	types.AppPortProfileScopeTenant,
	types.AppPortProfileScopeProvider,
	types.AppPortProfileScopeSystem,
}

var portRegex = regexp.MustCompile(types.AppPortProfilePortRegex)

func IsAppPortProfileURN(idOrName string) bool {
	return urn.IsAppPortProfile(idOrName)
}

func FindAppPortProfile(ctx context.Context, c cav.Client, idOrName, vdcGroupID string) (*itypes.ApiResponseAppPortProfile, error) {
	if IsAppPortProfileURN(idOrName) {
		ep := endpoints.GetAppPortProfile()

		resp, err := c.Do(
			ctx,
			ep,
			cav.WithPathParam(ep.PathParams[0], idOrName),
		)
		if err != nil {
			return nil, err
		}

		return resp.Result().(*itypes.ApiResponseAppPortProfile), nil
	}

	var found []itypes.ApiResponseAppPortProfile
	for _, scope := range appPortProfileScopes {
		ep := endpoints.ListAppPortProfile()

		resp, err := c.Do(
			ctx,
			ep,
			cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("name==%s;scope==%s;_context==%s", idOrName, scope, vdcGroupID)),
		)
		if err != nil {
			continue
		}

		list := resp.Result().(*itypes.ApiResponseListAppPortProfile)
		if len(list.Values) == 0 {
			continue
		}
		if len(list.Values) > 1 {
			return nil, pkgerrors.Newf("found multiple application port profiles with the same name %q", idOrName)
		}

		found = append(found, list.Values[0])
	}

	if len(found) == 0 {
		return nil, &pkgerrors.APIError{Operation: "GetAppPortProfile", StatusCode: 404, Message: fmt.Sprintf("application port profile %q not found", idOrName)}
	}

	return &found[0], nil
}

func ValidateAppPortProfileApplicationPorts(ports []types.ParamsAppPortProfilePort) error {
	for _, port := range ports {
		switch port.Protocol {
		case types.AppPortProfileProtocolICMPv4, types.AppPortProfileProtocolICMPv6:
			if len(port.DestinationPorts) != 0 {
				return pkgerrors.Newf("port must be empty for protocol %s", port.Protocol)
			}
		case types.AppPortProfileProtocolTCP, types.AppPortProfileProtocolUDP:
			if len(port.DestinationPorts) == 0 {
				return pkgerrors.Newf("port is required for protocol %s", port.Protocol)
			}
			for _, p := range port.DestinationPorts {
				if !portRegex.MatchString(p) {
					return pkgerrors.Newf("port %s is invalid", p)
				}
			}
		default:
			return pkgerrors.Newf("protocol must be one of %v", types.AppPortProfileProtocols)
		}
	}

	return nil
}

func ToApiAppPortProfilePorts(ports []types.ParamsAppPortProfilePort) []itypes.ApiAppPortProfilePort {
	out := make([]itypes.ApiAppPortProfilePort, 0, len(ports))
	for _, p := range ports {
		out = append(out, itypes.ApiAppPortProfilePort{
			Protocol:         p.Protocol,
			DestinationPorts: p.DestinationPorts,
		})
	}

	return out
}
