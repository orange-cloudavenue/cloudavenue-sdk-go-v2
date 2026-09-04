/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

const (
	NetworkContextProfileScopeSystem   = "SYSTEM"
	NetworkContextProfileScopeProvider = "PROVIDER"
	NetworkContextProfileScopeTenant   = "TENANT"

	NetworkContextProfileAttributeTypeAppID      = "APP_ID"
	NetworkContextProfileAttributeTypeDomainName = "DOMAIN_NAME"

	NetworkContextProfileSubAttributeTypeTLSVersion     = "TLS_VERSION"
	NetworkContextProfileSubAttributeTypeTLSCipherSuite = "TLS_CIPHER_SUITE"
	NetworkContextProfileSubAttributeTypeCIFSSMBVersion = "CIFS_SMB_VERSION"
)

// NetworkContextProfileTLSVersionValues lists all valid TLS_VERSION sub-attribute values.
var NetworkContextProfileTLSVersionValues = []string{
	"TLS_V10",
	"TLS_V11",
	"TLS_V12",
	"TLS_V13",
}

// NetworkContextProfileTLSCipherSuiteValues lists all valid TLS_CIPHER_SUITE sub-attribute values.
var NetworkContextProfileTLSCipherSuiteValues = []string{
	"TLS_DHE_RSA_WITH_AES_128_CBC_SHA",
	"TLS_DHE_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_DHE_RSA_WITH_AES_256_CBC_SHA",
	"TLS_DHE_RSA_WITH_AES_256_CBC_SHA256",
	"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384",
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	"TLS_RSA_WITH_AES_128_CBC_SHA",
	"TLS_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_RSA_WITH_AES_256_CBC_SHA",
	"TLS_RSA_WITH_AES_256_CBC_SHA256",
	"TLS_RSA_WITH_AES_256_GCM_SHA384",
}

// NetworkContextProfileCIFSSMBVersionValues lists all valid CIFS_SMB_VERSION sub-attribute values.
var NetworkContextProfileCIFSSMBVersionValues = []string{
	"SMB_V1",
	"SMB_V2",
	"SMB_V3",
}

// NetworkContextProfileAppIDValues lists all valid Layer 7 application identifiers (APP_ID values).
var NetworkContextProfileAppIDValues = []string{
	"360ANTIV", "ACTIVDIR", "AMQP", "AVAST", "AVG", "AVIRA", "BDEFNDER", "BLAST",
	"CA_CERT", "CIFS", "CLDAP", "CTRXCGP", "CTRXGOTO", "CTRXICA", "DCERPC", "DHCP",
	"DIAMETER", "DNS", "EPIC", "ESET", "FPROT", "FTP", "GITHUB", "HTTP", "HTTP2",
	"IMAP", "KASPRSKY", "KERBEROS", "LDAP", "MAXDB", "MCAFEE", "MSSQL", "MYSQL",
	"NFS", "NNTP", "NTBIOSNS", "NTP", "OCSP", "ORACLE", "PANDA", "PCOIP", "POP3",
	"RADIUS", "RDP", "RTCP", "RTP", "RTSP", "SIP", "SMTP", "SNMP", "SSH", "SSL",
	"SYMUPDAT", "SYSLOG", "TELNET", "TFTP", "VNC", "WINS",
}

type (
	ParamsListNetworkContextProfile struct {
		// VDCGroupID is the ID of the Vdc Group owning the Network Context Profiles to list.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the Network Context Profiles to list.
		VDCGroupName string
	}

	ParamsGetNetworkContextProfile struct {
		// ID is the unique identifier of the Network Context Profile.
		ID string

		// Name is the name of the Network Context Profile.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Network Context Profile. Required
		// when Name is used to look up the profile.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Network Context Profile.
		VDCGroupName string
	}

	ParamsCreateNetworkContextProfile struct {
		// Name is the name of the Network Context Profile.
		Name string

		// Description is the description of the Network Context Profile.
		Description string

		// VDCGroupID is the ID of the Vdc Group owning this Network Context Profile.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Network Context Profile.
		VDCGroupName string

		// Attributes is the list of Layer 7 attributes (APP_ID or DOMAIN_NAME) for this profile.
		Attributes []ParamsNetworkContextProfileAttribute
	}

	ParamsNetworkContextProfileAttribute struct {
		// Type is the attribute type (APP_ID or DOMAIN_NAME).
		Type string

		// Values is the list of values for this attribute.
		Values []string

		// SubAttributes provides optional refinements (only valid under APP_ID attributes; e.g.
		// TLS_VERSION, TLS_CIPHER_SUITE, CIFS_SMB_VERSION).
		SubAttributes []ParamsNetworkContextProfileSubAttribute
	}

	ParamsNetworkContextProfileSubAttribute struct {
		// Type identifies the sub-attribute (TLS_VERSION, TLS_CIPHER_SUITE, CIFS_SMB_VERSION).
		Type string

		// Values is the list of allowed values for this sub-attribute.
		Values []string
	}

	ParamsUpdateNetworkContextProfile struct {
		// ID is the unique identifier of the Network Context Profile to update.
		ID string

		// Name is the new name of the Network Context Profile.
		Name string

		// Description is the new description of the Network Context Profile.
		Description string

		// Attributes is the full list of Layer 7 attributes for this profile.
		Attributes []ParamsNetworkContextProfileAttribute
	}

	ParamsDeleteNetworkContextProfile struct {
		// ID is the unique identifier of the Network Context Profile to delete.
		ID string

		// Name is the name of the Network Context Profile to delete.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Network Context Profile.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Network Context Profile.
		VDCGroupName string
	}

	ParamsGetNetworkContextProfileAttributes struct {
		// VDCGroupID is the ID of the Vdc Group to fetch the live attribute catalog for.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group to fetch the live attribute catalog for.
		VDCGroupName string
	}

	// * List
	ModelListNetworkContextProfile struct {
		NetworkContextProfiles []ModelGetNetworkContextProfile `documentation:"List of Network Context Profiles"`
	}

	// * Get
	ModelGetNetworkContextProfile struct {
		ID          string `documentation:"ID of the Network Context Profile"`
		Name        string `documentation:"Name of the Network Context Profile"`
		Description string `documentation:"Description of the Network Context Profile"`
		Scope       string `documentation:"Scope of the Network Context Profile (SYSTEM, PROVIDER, TENANT)"`
		OrgID       string `documentation:"ID of the Org owning the Network Context Profile (TENANT scope only)"`

		Attributes []ModelNetworkContextProfileAttribute `documentation:"List of Layer 7 attributes for this profile"`
	}

	ModelNetworkContextProfileAttribute struct {
		Type          string                                   `documentation:"Type of the attribute (APP_ID or DOMAIN_NAME)"`
		Values        []string                                 `documentation:"List of values for this attribute"`
		SubAttributes []ModelNetworkContextProfileSubAttribute `documentation:"List of sub-attributes refining this attribute"`
	}

	ModelNetworkContextProfileSubAttribute struct {
		Type   string   `documentation:"Type of the sub-attribute (TLS_VERSION, TLS_CIPHER_SUITE, CIFS_SMB_VERSION)"`
		Values []string `documentation:"List of allowed values for this sub-attribute"`
	}

	// ModelGetNetworkContextProfileAttributesCatalog holds the live, server-computed catalog of
	// valid attribute values returned by GetNetworkContextProfileAttributes.
	ModelGetNetworkContextProfileAttributesCatalog struct {
		AppIDValues      []string `documentation:"List of valid APP_ID values on this platform"`
		DomainNameValues []string `documentation:"List of valid DOMAIN_NAME values on this platform"`
	}
)
