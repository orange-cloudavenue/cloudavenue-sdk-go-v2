/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// * Request / Response API (JSON)

// APIRequestLDAPTest represents the request body for TestLDAP.
type APIRequestLDAPTest struct {
	HostName                 string                         `json:"hostName,omitempty"`
	Port                     *int                           `json:"port,omitempty"`
	IsSSL                    *bool                          `json:"isSsl,omitempty"`
	IsSSLAcceptAll           *bool                          `json:"isSslAcceptAll,omitempty"`
	Realm                    string                         `json:"realm,omitempty"`
	PagedSearchDisabled      *bool                          `json:"pagedSearchDisabled,omitempty"`
	PageSize                 *int                           `json:"pageSize,omitempty"`
	MaxResults               *int                           `json:"maxResults,omitempty"`
	MaxUserGroups            *int                           `json:"maxUserGroups,omitempty"`
	SearchBase               string                         `json:"searchBase,omitempty"`
	UserName                 string                         `json:"userName,omitempty"`
	Password                 string                         `json:"password,omitempty"`
	AuthenticationMechanism  string                         `json:"authenticationMechanism,omitempty"`
	GroupSearchBase          string                         `json:"groupSearchBase,omitempty"`
	IsGroupSearchBaseEnabled *bool                          `json:"isGroupSearchBaseEnabled,omitempty"`
	ConnectorType            string                         `json:"connectorType,omitempty"`
	UserAttributes           *APIRequestLDAPUserAttributes  `json:"userAttributes,omitempty"`
	GroupAttributes          *APIRequestLDAPGroupAttributes `json:"groupAttributes,omitempty"`
	UseExternalKerberos      *bool                          `json:"useExternalKerberos,omitempty"`
	CustomUIButtonLabel      string                         `json:"customUiButtonLabel,omitempty"`
}

// APIRequestLDAPUserAttributes represents LDAP user attribute mapping for TestLDAP.
type APIRequestLDAPUserAttributes struct {
	ObjectClass               string `json:"objectClass,omitempty"`
	ObjectIdentifier          string `json:"objectIdentifier,omitempty"`
	UserName                  string `json:"userName,omitempty"`
	Email                     string `json:"email,omitempty"`
	FullName                  string `json:"fullName,omitempty"`
	GivenName                 string `json:"givenName,omitempty"`
	Surname                   string `json:"surname,omitempty"`
	Telephone                 string `json:"telephone,omitempty"`
	GroupMembershipIdentifier string `json:"groupMembershipIdentifier,omitempty"`
	GroupBackLinkIdentifier   string `json:"groupBackLinkIdentifier,omitempty"`
}

// APIRequestLDAPGroupAttributes represents LDAP group attribute mapping for TestLDAP.
type APIRequestLDAPGroupAttributes struct {
	ObjectClass          string `json:"objectClass,omitempty"`
	ObjectIdentifier     string `json:"objectIdentifier,omitempty"`
	GroupName            string `json:"groupName,omitempty"`
	Membership           string `json:"membership,omitempty"`
	MembershipIdentifier string `json:"membershipIdentifier,omitempty"`
	BackLinkIdentifier   string `json:"backLinkIdentifier,omitempty"`
}

// APIResponseLDAPTestResult represents the response of TestLDAP.
type APIResponseLDAPTestResult struct {
	ConnectionTest *APIResponseLDAPConnectionTest    `json:"connectionTest,omitempty"`
	SettingsTest   []APIResponseLDAPSettingsTestItem `json:"settingsTest,omitempty"`
}

// APIResponseLDAPConnectionTest represents LDAP connection test result.
type APIResponseLDAPConnectionTest struct {
	Successful *bool                 `json:"successful,omitempty"`
	Error      *APIResponseLDAPError `json:"error,omitempty"`
}

// APIResponseLDAPError represents LDAP test error details.
type APIResponseLDAPError struct {
	MinorErrorCode string `json:"minorErrorCode,omitempty"`
	Message        string `json:"message,omitempty"`
	StackTrace     string `json:"stackTrace,omitempty"`
}

// APIResponseLDAPSettingsTestItem represents single LDAP settings test result.
type APIResponseLDAPSettingsTestItem struct {
	Attribute      string `json:"attribute,omitempty"`
	AttributeValue string `json:"attributeValue,omitempty"`
	Result         string `json:"result,omitempty"`
	Successful     *bool  `json:"successful,omitempty"`
}

// APIResponseLDAPUser represents a single LDAP user returned by SearchLDAPUsers.
type APIResponseLDAPUser struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}

// APIResponseLDAPGroup represents a single LDAP group returned by SearchLDAPGroups.
type APIResponseLDAPGroup struct {
	Name string `json:"name,omitempty"`
}

// ToModel converts the APIResponseLDAPTestResult to ModelLDAPTestResult.
func (api *APIResponseLDAPTestResult) ToModel() *types.ModelLDAPTestResult {
	if api == nil {
		return nil
	}

	return &types.ModelLDAPTestResult{
		Success: api.ConnectionTest != nil && api.ConnectionTest.Successful != nil && *api.ConnectionTest.Successful,
		Message: ldapTestMessage(api),
	}
}

// ToModels converts each LDAP user to ModelLDAPUser.
func LDAPUsersToModel(api []APIResponseLDAPUser) []*types.ModelLDAPUser {
	result := make([]*types.ModelLDAPUser, len(api))
	for i, u := range api {
		result[i] = u.ToModel()
	}

	return result
}

// ToModel converts a single LDAP user to ModelLDAPUser.
func (api *APIResponseLDAPUser) ToModel() *types.ModelLDAPUser {
	if api == nil {
		return nil
	}

	return &types.ModelLDAPUser{
		Name:     api.Username,
		Email:    api.Email,
		FullName: api.Fullname,
	}
}

// ToModels converts each LDAP group to ModelLDAPGroup.
func LDAPGroupsToModel(api []APIResponseLDAPGroup) []*types.ModelLDAPGroup {
	result := make([]*types.ModelLDAPGroup, len(api))
	for i, g := range api {
		result[i] = g.ToModel()
	}

	return result
}

// ToModel converts a single LDAP group to ModelLDAPGroup.
func (api *APIResponseLDAPGroup) ToModel() *types.ModelLDAPGroup {
	if api == nil {
		return nil
	}

	return &types.ModelLDAPGroup{
		Name: api.Name,
	}
}

// LDAPConfigToAPIRequest converts public LDAP params to internal API request config.
// It fails when params carry fields the LDAP test payload cannot express, so no
// caller can silently drop them.
func LDAPConfigToAPIRequest(p types.ParamsTestLDAP) (APIRequestLDAPTest, error) {
	if err := p.Validate(); err != nil {
		return APIRequestLDAPTest{}, fmt.Errorf("ldap config: %w", err)
	}

	var groupSearchBaseEnabled *bool
	if p.GroupSearchBase != "" {
		enabled := true
		groupSearchBaseEnabled = &enabled
	}

	return APIRequestLDAPTest{
		HostName:                 p.Host,
		Port:                     p.Port,
		IsSSL:                    p.SSLEnabled,
		SearchBase:               p.BaseDN,
		UserName:                 p.BindUser,
		Password:                 p.BindPassword,
		AuthenticationMechanism:  p.AuthenticationMechanism,
		GroupSearchBase:          p.GroupSearchBase,
		IsGroupSearchBaseEnabled: groupSearchBaseEnabled,
		UserAttributes: &APIRequestLDAPUserAttributes{
			ObjectClass: p.UserObjectClass,
			UserName:    p.UserNameAttribute,
		},
		GroupAttributes: &APIRequestLDAPGroupAttributes{
			ObjectClass: p.GroupObjectClass,
			GroupName:   p.GroupNameAttribute,
			Membership:  p.GroupMemberAttribute,
		},
	}, nil
}

func ldapTestMessage(api *APIResponseLDAPTestResult) string {
	if api == nil {
		return ""
	}
	if api.ConnectionTest != nil && api.ConnectionTest.Error != nil && api.ConnectionTest.Error.Message != "" {
		return api.ConnectionTest.Error.Message
	}
	for _, item := range api.SettingsTest {
		if item.Result != "" {
			return item.Result
		}
	}
	return ""
}
