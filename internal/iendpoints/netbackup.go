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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path netbackup.go -output netbackup

func init() {
	// GetNetbackupToken
	cav.Endpoint{
		DocumentationURL: "https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/auth/token",
		Name:             "GetNetbackupToken",
		Description:      "Get NetBackup OAuth2 token using password grant",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendNetBackup,
		PathTemplate:     pathNetBackupAuthToken,
		BodyRequestType:  netbackupTokenRequest{},
		ResponseType:     netbackupTokenResponse{},
	}.Register()

	// RefreshNetbackupToken
	cav.Endpoint{
		DocumentationURL: "https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/auth/token",
		Name:             "RefreshNetbackupToken",
		Description:      "Refresh NetBackup OAuth2 token using refresh token grant",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendNetBackup,
		PathTemplate:     pathNetBackupAuthToken,
		BodyRequestType:  netbackupRefreshTokenRequest{},
		ResponseType:     netbackupTokenResponse{},
	}.Register()

	// ListNetbackupInventory
	cav.Endpoint{
		DocumentationURL: "https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/inventory",
		Name:             "ListNetbackupInventory",
		Description:      "List NetBackup inventory",
		Method:           cav.MethodGET,
		Backend:          cav.BackendNetBackup,
		PathTemplate:     pathNetBackupInventory,
		ResponseType:     itypes.APIResponseListNetbackupInventory{},
	}.Register()

	// ListNetbackupMachines
	cav.Endpoint{
		DocumentationURL: "https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/machines",
		Name:             "ListNetbackupMachines",
		Description:      "List NetBackup machines",
		Method:           cav.MethodGET,
		Backend:          cav.BackendNetBackup,
		PathTemplate:     pathNetBackupMachines,
		ResponseType:     itypes.APIResponseListNetbackupMachines{},
	}.Register()

	// GetNetbackupProtectionLevel
	cav.Endpoint{
		DocumentationURL: "https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/protection-levels/{id}",
		Name:             "GetNetbackupProtectionLevel",
		Description:      "Get NetBackup protection level by ID",
		Method:           cav.MethodGET,
		Backend:          cav.BackendNetBackup,
		PathTemplate:     pathNetBackupProtectionLevelByID,
		PathParams: []cav.PathParam{
			{
				Name:        "id",
				Description: "The protection level ID",
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseNetbackupProtectionLevel{},
	}.Register()

	// ListNetbackupProtectionLevels
	cav.Endpoint{
		DocumentationURL: "https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/protection-levels",
		Name:             "ListNetbackupProtectionLevels",
		Description:      "List NetBackup protection levels",
		Method:           cav.MethodGET,
		Backend:          cav.BackendNetBackup,
		PathTemplate:     pathNetBackupProtectionLevels,
		ResponseType:     itypes.APIResponseListNetbackupProtectionLevels{},
	}.Register()

	// ProtectNetbackupMachine
	cav.Endpoint{
		DocumentationURL: "https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/machines/{id}/protect",
		Name:             "ProtectNetbackupMachine",
		Description:      "Protect a NetBackup machine",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendNetBackup,
		PathTemplate:     pathNetBackupMachineProtect,
		PathParams: []cav.PathParam{
			{
				Name:        "id",
				Description: "The machine ID",
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestNetbackupProtectMachine{},
		ResponseType:    itypes.APIResponseNetbackupProtectMachine{},
	}.Register()
}

type netbackupTokenRequest struct {
	GrantType string `json:"grant_type"` //nolint:tagliatelle
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type netbackupRefreshTokenRequest struct {
	GrantType    string `json:"grant_type"`    //nolint:tagliatelle
	RefreshToken string `json:"refresh_token"` //nolint:tagliatelle
}

type netbackupTokenResponse struct {
	AccessToken  string `json:"access_token"`  //nolint:tagliatelle
	TokenType    string `json:"token_type"`    //nolint:tagliatelle
	ExpiresIn    int    `json:"expires_in"`    //nolint:tagliatelle
	RefreshToken string `json:"refresh_token"` //nolint:tagliatelle //nolint:tagliatelle
}
