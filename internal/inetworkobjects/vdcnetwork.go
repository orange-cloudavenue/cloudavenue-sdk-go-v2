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
	"strconv"
	"strings"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// VDCNetworkPageSize is the page size for Org VDC Network listings. VCD
// enforces this maximum for the resource.
const VDCNetworkPageSize = 32

// ListVDCNetworkModels lists the Org VDC Networks matching a cloudapi filter,
// fetching every page, and returns models carrying a legacy network HREF.
func ListVDCNetworkModels(ctx context.Context, c cav.Client, filter, baseURL string) ([]types.ModelGetVDCNetwork, error) {
	ep := endpoints.ListVDCNetwork()
	networks := make([]types.ModelGetVDCNetwork, 0)

	for page := 1; ; page++ {
		resp, err := c.Do(
			ctx,
			ep,
			cav.WithQueryParam(ep.QueryParams[0], filter),
			cav.WithQueryParam(ep.QueryParams[1], strconv.Itoa(VDCNetworkPageSize)),
			cav.WithQueryParam(ep.QueryParams[2], strconv.Itoa(page)),
		)
		if err != nil {
			return nil, err
		}

		list, ok := resp.Result().(*itypes.APIResponseListVDCNetwork)
		if !ok || list == nil {
			return nil, fmt.Errorf("unexpected list response type %T", resp.Result())
		}

		pageModel := list.ToModel()
		for i := range pageModel.VDCNetworks {
			LegacyVDCNetworkHREF(&pageModel.VDCNetworks[i], baseURL)
		}
		networks = append(networks, pageModel.VDCNetworks...)

		if len(list.Values) == 0 ||
			(list.PageCount > 0 && page >= list.PageCount) ||
			(list.PageCount == 0 && len(list.Values) < VDCNetworkPageSize) {
			break
		}
	}

	return networks, nil
}

// LegacyVDCNetworkHREF fills the legacy network HREF of a model built from a
// cloudapi payload, which carries no href, while the legacy vApp runtime
// network contract (Configuration.ParentNetwork) is keyed on it.
//
// The live VCD answers GET {base}/api/network/{uuid} with an OrgVdcNetwork whose
// own href attribute is that same URL, so it is the reference consumers must
// send; the /api/admin/network variant is admin-scoped and rejected with
// organization credentials.
func LegacyVDCNetworkHREF(network *types.ModelGetVDCNetwork, baseURL string) {
	if network == nil || network.HREF != "" || baseURL == "" || network.ID == "" {
		return
	}

	id := network.ID
	if urn.IsNetwork(id) {
		extracted, err := extractor.ExtractUUID(id)
		if err != nil {
			return
		}
		id = extracted
	}
	if id == "" {
		return
	}

	network.HREF = strings.TrimRight(baseURL, "/") + "/api/network/" + id
}
