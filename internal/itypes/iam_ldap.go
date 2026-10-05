/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

// * Request / Response API (JSON)

// APIRequestLDAPConfig represents the LDAP configuration used in test/sync requests.
type APIRequestLDAPConfig struct {
	Host                    string `json:"host,omitempty"`
	Port                    *int   `json:"port,omitempty"`
	BindUser                string `json:"bindUser,omitempty"`
	BindPassword            string `json:"bindPassword,omitempty"`
	BaseDN                  string `json:"baseDn,omitempty"`
	UserSearchBase          string `json:"userSearchBase,omitempty"`
	GroupSearchBase         string `json:"groupSearchBase,omitempty"`
	UserObjectClass         string `json:"userObjectClass,omitempty"`
	GroupObjectClass        string `json:"groupObjectClass,omitempty"`
	UserNameAttribute       string `json:"userNameAttribute,omitempty"`
	GroupNameAttribute      string `json:"groupNameAttribute,omitempty"`
	GroupMemberAttribute    string `json:"groupMemberAttribute,omitempty"`
	AuthenticationMechanism string `json:"authenticationMechanism,omitempty"`
	SSLEnabled              *bool  `json:"sslEnabled,omitempty"`
	SSLTrustCertificate     string `json:"sslTrustCertificate,omitempty"`
	ConnectionTimeout       *int   `json:"connectionTimeout,omitempty"`
	ReadTimeout             *int   `json:"readTimeout,omitempty"`
	Enabled                 *bool  `json:"enabled,omitempty"`
}

// APIRequestLDAPTest represents the request body for TestLDAP.
type APIRequestLDAPTest struct {
	Config APIRequestLDAPConfig `json:"config,omitempty"`
}

// APIResponseLDAPTestResult represents the response of TestLDAP.
type APIResponseLDAPTestResult struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// APIResponseLDAPUser represents a single LDAP user returned by SearchLDAPUsers.
type APIResponseLDAPUser struct {
	Name     string `json:"name,omitempty"`
	DN       string `json:"dn,omitempty"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"fullName,omitempty"`
}

// APIResponseLDAPUsers represents the wrapper for a list of LDAP users.
type APIResponseLDAPUsers struct {
	Users []APIResponseLDAPUser `json:"users,omitempty"`
}

// APIResponseLDAPGroup represents a single LDAP group returned by SearchLDAPGroups.
type APIResponseLDAPGroup struct {
	Name    string   `json:"name,omitempty"`
	DN      string   `json:"dn,omitempty"`
	Members []string `json:"members,omitempty"`
}

// APIResponseLDAPGroups represents the wrapper for a list of LDAP groups.
type APIResponseLDAPGroups struct {
	Groups []APIResponseLDAPGroup `json:"groups,omitempty"`
}

// ToModel converts the APIResponseLDAPTestResult to ModelLDAPTestResult.
func (api *APIResponseLDAPTestResult) ToModel() *types.ModelLDAPTestResult {
	if api == nil {
		return nil
	}

	return &types.ModelLDAPTestResult{
		Success: api.Success,
		Message: api.Message,
	}
}

// ToModel converts each LDAP user to ModelLDAPUser.
func (api *APIResponseLDAPUsers) ToModel() []*types.ModelLDAPUser {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelLDAPUser, len(api.Users))
	for i, u := range api.Users {
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
		Name:     api.Name,
		DN:       api.DN,
		Email:    api.Email,
		FullName: api.FullName,
	}
}

// ToModel converts each LDAP group to ModelLDAPGroup.
func (api *APIResponseLDAPGroups) ToModel() []*types.ModelLDAPGroup {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelLDAPGroup, len(api.Groups))
	for i, g := range api.Groups {
		result[i] = g.ToModel()
	}

	return result
}

// ToModel converts a single LDAP group to ModelLDAPGroup.
func (api *APIResponseLDAPGroup) ToModel() *types.ModelLDAPGroup {
	if api == nil {
		return nil
	}

	members := make([]string, len(api.Members))
	copy(members, api.Members)

	return &types.ModelLDAPGroup{
		Name:    api.Name,
		DN:      api.DN,
		Members: members,
	}
}

// LDAPConfigToAPIRequest converts public LDAP params to internal API request config.
func LDAPConfigToAPIRequest(p types.ParamsTestLDAP) APIRequestLDAPConfig {
	return APIRequestLDAPConfig{
		Host:                    p.Host,
		Port:                    p.Port,
		BindUser:                p.BindUser,
		BindPassword:            p.BindPassword,
		BaseDN:                  p.BaseDN,
		UserSearchBase:          p.UserSearchBase,
		GroupSearchBase:         p.GroupSearchBase,
		UserObjectClass:         p.UserObjectClass,
		GroupObjectClass:        p.GroupObjectClass,
		UserNameAttribute:       p.UserNameAttribute,
		GroupNameAttribute:      p.GroupNameAttribute,
		GroupMemberAttribute:    p.GroupMemberAttribute,
		AuthenticationMechanism: p.AuthenticationMechanism,
		SSLEnabled:              p.SSLEnabled,
		SSLTrustCertificate:     p.SSLTrustCertificate,
		ConnectionTimeout:       p.ConnectionTimeout,
		ReadTimeout:             p.ReadTimeout,
		Enabled:                 p.Enabled,
	}
}
