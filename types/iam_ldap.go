/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

import (
	"errors"
	"fmt"
	"strings"
)

// ErrUnsupportedLDAPTestParams is returned when a ParamsTestLDAP field has no
// counterpart in the LDAP test request payload.
var ErrUnsupportedLDAPTestParams = errors.New("unsupported params for endpoint")

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

	// Email is the LDAP user email address.
	Email string `documentation:"Email address of the LDAP user"`

	// FullName is the LDAP user full name.
	FullName string `documentation:"Full name of the LDAP user"`
}

// ModelLDAPGroup represents a group discovered in LDAP.
type ModelLDAPGroup struct {
	// Name is the LDAP group name.
	Name string `documentation:"Name of the LDAP group"`
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
	// Not supported by the LDAP test endpoint: Validate rejects it when set.
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
	// Not supported by the LDAP test endpoint: Validate rejects it when set.
	SSLTrustCertificate string

	// ConnectionTimeout is the connection timeout in seconds.
	// Not supported by the LDAP test endpoint: Validate rejects it when set.
	ConnectionTimeout *int

	// ReadTimeout is the read timeout in seconds.
	// Not supported by the LDAP test endpoint: Validate rejects it when set.
	ReadTimeout *int

	// Enabled indicates whether LDAP authentication is enabled.
	// Not supported by the LDAP test endpoint: Validate rejects it when non-nil.
	Enabled *bool
}

// Validate checks ParamsTestLDAP structural constraints.
func (p ParamsTestLDAP) Validate() error {
	if p.Host == "" {
		return fmt.Errorf("host is required")
	}

	if names := p.unsupportedFields(); len(names) > 0 {
		return fmt.Errorf("%w: %s", ErrUnsupportedLDAPTestParams, strings.Join(names, ", "))
	}

	return nil
}

// unsupportedFields returns the lowerCamelCase names of the fields that are set
// but absent from the LDAP test request payload.
func (p ParamsTestLDAP) unsupportedFields() []string {
	var names []string

	if p.UserSearchBase != "" {
		names = append(names, "userSearchBase")
	}
	if p.SSLTrustCertificate != "" {
		names = append(names, "sslTrustCertificate")
	}
	if p.ConnectionTimeout != nil {
		names = append(names, "connectionTimeout")
	}
	if p.ReadTimeout != nil {
		names = append(names, "readTimeout")
	}
	if p.Enabled != nil {
		names = append(names, "enabled")
	}

	return names
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
}

// Validate checks ParamsSearchLDAP structural constraints.
func (p ParamsSearchLDAP) Validate() error {
	return nil
}
