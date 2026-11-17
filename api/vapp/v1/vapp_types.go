/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

type getVAppByIDParams struct {
	ID string
}

type composeVAppParams struct {
	VDCID string
	Body  itypes.APIRequestComposeVApp
}

type instantiateVAppTemplateParams struct {
	VDCID string
	Body  itypes.APIRequestInstantiateVAppTemplate
}

type updateVAppByIDParams struct {
	ID   string
	Body itypes.APIRequestUpdateVApp
}

type deleteVAppByIDParams struct {
	ID string
}

type updateVAppNetworkConfigSectionParams struct {
	ID   string
	Body itypes.APIResponseVAppNetworkConfigSection
}

type updateVAppGuestPropertiesParams struct {
	ID   string
	Body itypes.APIResponseGuestProperties
}

type setVAppACLParams struct {
	ID   string
	Body itypes.ControlAccessParams
}

type undeployVAppParams struct {
	ID   string
	Body itypes.APIRequestUndeployVApp
}

type listVAppParams struct {
	VDCID string
}
