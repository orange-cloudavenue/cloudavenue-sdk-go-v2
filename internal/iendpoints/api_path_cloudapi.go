/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

const (
	// pathCloudAPIOrgs lists organizations in VMware Cloud Director CloudAPI.
	pathCloudAPIOrgs = "/cloudapi/1.0.0/orgs"
	// pathCloudAPIUsers lists and creates VMware Cloud Director CloudAPI users.
	pathCloudAPIUsers = "/cloudapi/1.0.0/users"
	// pathCloudAPIUser addresses one CloudAPI user by URN.
	pathCloudAPIUser = "/cloudapi/1.0.0/users/{userUrn}"
	// pathCloudAPIUserChangePassword changes one CloudAPI user's password.
	pathCloudAPIUserChangePassword = "/cloudapi/1.0.0/users/{userUrn}/changePassword"
	// pathCloudAPIUserTakeOwnership transfers one user's owned entities.
	pathCloudAPIUserTakeOwnership = "/cloudapi/1.0.0/users/{userUrn}/takeOwnership"
	// pathApplicationPortProfilesBase lists and creates application port profiles.
	pathApplicationPortProfilesBase = "/cloudapi/1.0.0/applicationPortProfiles"
	// pathApplicationPortProfiles addresses one application port profile by ID.
	pathApplicationPortProfiles = "/cloudapi/1.0.0/applicationPortProfiles/{appPortProfileId}"
	// pathCertificateLibraryBase lists and creates certificate library items.
	pathCertificateLibraryBase = "/cloudapi/1.0.0/ssl/certificateLibrary"
	// pathCertificateLibrary addresses one certificate library item by ID.
	pathCertificateLibrary = "/cloudapi/1.0.0/ssl/certificateLibrary/{id}"
	// pathCertificateLibraryConsumers lists or mutates certificate consumers.
	pathCertificateLibraryConsumers = "/cloudapi/1.0.0/ssl/certificateLibrary/{certLibraryItemId}/consumers"
	// pathFirewallGroupSummaries lists firewall group summaries.
	pathFirewallGroupSummaries = "/cloudapi/1.0.0/firewallGroups/summaries"
	// pathFirewallGroupsBase lists and creates firewall groups.
	pathFirewallGroupsBase = "/cloudapi/1.0.0/firewallGroups"
	// pathFirewallGroups addresses one firewall group by ID.
	pathFirewallGroups = "/cloudapi/1.0.0/firewallGroups/{firewallGroupId}"
	// pathNetworkContextProfilesBase lists and creates network context profiles.
	pathNetworkContextProfilesBase = "/cloudapi/1.0.0/networkContextProfiles"
	// pathNetworkContextProfiles addresses one network context profile by ID.
	pathNetworkContextProfiles = "/cloudapi/1.0.0/networkContextProfiles/{networkContextProfileId}"
	// pathNetworkContextProfileAttributes lists available network context profile attributes.
	pathNetworkContextProfileAttributes = "/cloudapi/1.0.0/networkContextProfiles/attributes"
	// pathTrustedCertificatesBase lists and creates trusted certificates.
	pathTrustedCertificatesBase = "/cloudapi/1.0.0/ssl/trustedCertificates"
	// pathTrustedCertificates addresses one trusted certificate by ID.
	pathTrustedCertificates = "/cloudapi/1.0.0/ssl/trustedCertificates/{trustedCertificate}"
	// pathOrgVDCNetworksBase lists and creates Org VDC networks.
	pathOrgVDCNetworksBase = "/cloudapi/1.0.0/orgVdcNetworks"
	// pathOrgVDCNetworks addresses one Org VDC network by ID.
	pathOrgVDCNetworks = "/cloudapi/1.0.0/orgVdcNetworks/{vdcNetworkId}"
	// pathCatalogAccessControl manages catalog access controls by catalog URN.
	pathCatalogAccessControl = "/cloudapi/1.0.0/catalogs/{catalogUrn}/accessControls"
	// pathOrgVDCNetworkDHCP addresses DHCP configuration for one Org VDC network.
	pathOrgVDCNetworkDHCP = "/cloudapi/1.0.0/orgVdcNetworks/{vdcNetworkId}/dhcp"
	// pathVDCGroups lists and creates VDC groups.
	pathVDCGroups = "/cloudapi/1.0.0/vdcGroups"
	// pathVDCGroupByID addresses one VDC group by ID.
	pathVDCGroupByID = "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}"
	// pathDFWPolicies manages distributed firewall policies for a VDC group.
	pathDFWPolicies = "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}/dfwPolicies"
	// pathDFWPolicyDefault addresses the default distributed firewall policy.
	pathDFWPolicyDefault = "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}/dfwPolicies/default"
	// pathDFWPolicyDefaultRules manages rules of the default distributed firewall policy.
	pathDFWPolicyDefaultRules = "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}/dfwPolicies/default/rules"
	// pathTokens manages API tokens.
	// pathTokenByID gets, updates, or deletes one token by ID.
	// pathTokenGet lists tokens through the compatibility get route.
	pathTokenGet = "/cloudapi/1.0.0/tokens/get/"
	// pathTokenGetByID gets one token through the compatibility get route.
	pathTokenGetByID = "/cloudapi/1.0.0/tokens/id/get/"
	// pathTokenCreate creates a token through the compatibility post route.
	pathTokenCreate = "/cloudapi/1.0.0/tokens/post/"
	// pathTokenUpdateByID updates one token through the compatibility put route.
	pathTokenUpdateByID = "/cloudapi/1.0.0/tokens/id/put/"
	// pathTokenDeleteByID deletes one token through the compatibility delete route.
	pathTokenDeleteByID = "/cloudapi/1.0.0/tokens/id/delete/"
	// pathLDAP manages LDAP operations.
	// pathLDAPTest tests LDAP configuration.
	pathLDAPTest = "/cloudapi/1.0.0/ldap/{orgId}/test"
	// pathLDAPSync triggers LDAP synchronization.
	pathLDAPSync = "/cloudapi/1.0.0/ldap/{orgId}/sync"
	// pathLDAPSearchUsers searches LDAP users.
	pathLDAPSearchUsers = "/cloudapi/1.0.0/ldap/{orgId}/search/user"
	// pathLDAPSearchGroups searches LDAP groups.
	pathLDAPSearchGroups = "/cloudapi/1.0.0/ldap/{orgId}/search/group"
	// pathGlobalRoles manages global roles.
	pathGlobalRoles = "/cloudapi/1.0.0/globalRoles"
	// pathGlobalRoleByID addresses one global role by ID.
	pathGlobalRoleByID = "/cloudapi/1.0.0/globalRoles/{id}"
	// pathGlobalRoleRights manages rights for one global role.
	pathGlobalRoleRights = "/cloudapi/1.0.0/globalRoles/{id}/rights"
	// pathGlobalRoleTenants manages tenant publication for one global role.
	pathGlobalRoleTenants = "/cloudapi/1.0.0/globalRoles/{id}/tenants"
	// pathGlobalRoleTenantsPublish publishes a global role to selected tenants.
	pathGlobalRoleTenantsPublish = "/cloudapi/1.0.0/globalRoles/{id}/tenants/publish"
	// pathGlobalRoleTenantsUnpublish unpublishes a global role from selected tenants.
	pathGlobalRoleTenantsUnpublish = "/cloudapi/1.0.0/globalRoles/{id}/tenants/unpublish"
	// pathGlobalRoleTenantsPublishAll publishes a global role to all tenants.
	pathGlobalRoleTenantsPublishAll = "/cloudapi/1.0.0/globalRoles/{id}/tenants/publishAll"
	// pathGlobalRoleTenantsUnpublishAll unpublishes a global role from all tenants.
	pathGlobalRoleTenantsUnpublishAll = "/cloudapi/1.0.0/globalRoles/{id}/tenants/unpublishAll"
	// pathOrgVDCStoragePolicies lists Org VDC storage policies.
	pathOrgVDCStoragePolicies = "/cloudapi/1.0.0/orgVdcStoragePolicies"
)

const (
	pathParamAppPortProfileID        = "appPortProfileId"
	pathParamCertLibraryItemID       = "certLibraryItemId"
	pathParamFirewallGroupID         = "firewallGroupId"
	pathParamNetworkContextProfileID = "networkContextProfileId"
	pathParamTrustedCertificate      = "trustedCertificate"
	pathParamVDCGroupID              = "vdcGroupId"
	pathParamVDCNetworkID            = "vdcNetworkId"
	pathParamCatalogUrn              = "catalogUrn"
	pathParamUserUrn                 = "userUrn"
)
