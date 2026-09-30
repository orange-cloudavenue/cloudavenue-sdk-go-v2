/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

import "fmt"

// ModelLDAPTestResult represents the result of an LDAP connection test.
type ModelLDAPTestResult struct {
	// Success indicates whether the LDAP connection test succeeded.
	Success bool `documentation:"Indicates whether the LDAP connection test succeeded"`

	// Message provides details about the test result.
	Message string `documentation:"Message providing details about the test result"`
}

// ModelLDAPUser represents a user discovered in LDAP.
type ModelLDAPUser struct {
	// Name is the LDAP user name (sAMAccountName / uid).
	Name string `documentation:"Name of the LDAP user"`

	// DN is the distinguished name of the user.
	DN string `documentation:"Distinguished name of the LDAP user"`

	// Email is the LDAP user email address.
	Email string `documentation:"Email address of the LDAP user"`

	// FullName is the LDAP user full name.
	FullName string `documentation:"Full name of the LDAP user"`
}

// ModelLDAPGroup represents a group discovered in LDAP.
type ModelLDAPGroup struct {
	// Name is the LDAP group name.
	Name string `documentation:"Name of the LDAP group"`

	// DN is the distinguished name of the group.
	DN string `documentation:"Distinguished name of the LDAP group"`

	// Members is the list of distinguished names of the group members.
	Members []string `documentation:"List of distinguished names of the group members"`
}

// ParamsTestLDAP defines parameters for testing an LDAP connection.
type ParamsTestLDAP struct {
	// Host is the LDAP server hostname or IP address.
	Host string

	// Port is the LDAP server port.
	Port *int

	// BindUser is the DN of the account used to bind to LDAP.
	BindUser string

	// BindPassword is the password of the bind account.
	BindPassword string

	// BaseDN is the base DN used to root the LDAP directory search.
	BaseDN string

	// UserSearchBase is the base DN used to search for users.
	UserSearchBase string

	// GroupSearchBase is the base DN used to search for groups.
	GroupSearchBase string

	// UserObjectClass is the LDAP object class used to identify users.
	UserObjectClass string

	// GroupObjectClass is the LDAP object class used to identify groups.
	GroupObjectClass string

	// UserNameAttribute is the LDAP attribute holding the user name.
	UserNameAttribute string

	// GroupNameAttribute is the LDAP attribute holding the group name.
	GroupNameAttribute string

	// GroupMemberAttribute is the LDAP attribute holding the group members.
	GroupMemberAttribute string

	// AuthenticationMechanism is the SASL authentication mechanism.
	AuthenticationMechanism string

	// SSLEnabled indicates whether SSL/TLS is enabled.
	SSLEnabled *bool

	// SSLTrustCertificate is the trusted certificate (PEM) for SSL.
	SSLTrustCertificate string

	// ConnectionTimeout is the connection timeout in seconds.
	ConnectionTimeout *int

	// ReadTimeout is the read timeout in seconds.
	ReadTimeout *int

	// Enabled indicates whether LDAP authentication is enabled.
	Enabled *bool
}

// Validate checks ParamsTestLDAP structural constraints.
func (p ParamsTestLDAP) Validate() error {
	if p.Host == "" {
		return fmt.Errorf("host is required")
	}
	return nil
}

// ParamsSyncLDAP defines parameters for triggering an LDAP synchronization.
type ParamsSyncLDAP struct{}

// Validate checks ParamsSyncLDAP structural constraints.
func (p ParamsSyncLDAP) Validate() error {
	return nil
}

// ParamsSearchLDAP defines parameters for searching LDAP users or groups.
type ParamsSearchLDAP struct {
	// Filter is an LDAP search filter.
	Filter string

	// MaxResults is the maximum number of results to return.
	MaxResults string

	// PageSize is the page size for paginated searches.
	PageSize string

	// Page is the page index to fetch.
	Page string
}

// Validate checks ParamsSearchLDAP structural constraints.
func (p ParamsSearchLDAP) Validate() error {
	return nil
}
