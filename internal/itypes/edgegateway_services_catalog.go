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

var ListOfServices = []types.ModelCloudavenueServicesCatalog{
	{
		Category: "administration",
		Network:  "57.199.209.192/27",
		Services: []types.ModelCloudavenueServicesCatalogService{
			{
				Name:        "linux-repository",
				Description: "Linux (Debian, Ubuntu, CentOS) package repository",
				IPs:         []string{"57.199.209.214"},
				FQDNs:       []string{"repo.service.cav"},
				Ports:       []types.ModelCloudavenueServicesCatalogServicePort{{Port: 3142, Protocol: "tcp"}},
			},
			{
				Name:        "rhui-repository",
				Description: "Red Hat (RHUI) package repository",
				IPs:         []string{"57.199.209.197"},
				FQDNs:       []string{"rhui.service.cav"},
				Ports:       []types.ModelCloudavenueServicesCatalogServicePort{{Port: 8080, Protocol: "tcp"}},
			},
			{
				Name:        "windows-repository",
				Description: "Windows (WSUS) package repository",
				IPs:         []string{"57.199.209.212"},
				FQDNs:       []string{"wsus.service.cav"},
				Ports:       []types.ModelCloudavenueServicesCatalogServicePort{{Port: 8530, Protocol: "tcp"}, {Port: 8531, Protocol: "tcp"}},
			},
			{
				Name:        "windows-kms",
				Description: "Windows (KMS) license server",
				IPs:         []string{"57.199.209.210"},
				FQDNs:       []string{"kms.service.cav"},
				Ports:       []types.ModelCloudavenueServicesCatalogServicePort{{Port: 1688, Protocol: "tcp"}},
			},
			{
				Name:        "ntp",
				Description: "Network Time Protocol (NTP) server",
				IPs:         []string{"57.199.209.217", "57.199.209.218"},
				FQDNs:       []string{"ntp1.service.cav", "ntp2.service.cav"},
				Ports:       []types.ModelCloudavenueServicesCatalogServicePort{{Port: 123, Protocol: "udp"}},
			},
			{
				Name:             "dns-authoritative",
				Description:      "DNS authoritative server. Use for resolving cloudavenue services names",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/service-zone-dns/",
				IPs:              []string{"57.199.209.207", "57.199.209.208"},
				Ports:            []types.ModelCloudavenueServicesCatalogServicePort{{Port: 53, Protocol: "tcp"}, {Port: 53, Protocol: "udp"}},
			},
			{
				Name:             "dns-resolver",
				Description:      "DNS resolver. Use for resolving cloudavenue services names and public names",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/service-zone-dns/",
				IPs:              []string{"57.199.209.220", "57.199.209.221"},
				Ports:            []types.ModelCloudavenueServicesCatalogServicePort{{Port: 53, Protocol: "tcp"}, {Port: 53, Protocol: "udp"}},
			},
			{
				Name:             "smtp",
				Description:      "SMTP relay. Use for sending emails",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/smtp-service-2/",
				IPs:              []string{"57.199.209.206"},
				FQDNs:            []string{"smtp.service.cav"},
				Ports:            []types.ModelCloudavenueServicesCatalogServicePort{{Port: 25, Protocol: "tcp"}},
			},
		},
	},
	{
		Category: "s3",
		Network:  "194.206.55.5/32",
		Services: []types.ModelCloudavenueServicesCatalogService{{
			Name:             "s3-internal",
			Description:      "S3 internal service. Use for accessing S3 directly from the organization",
			DocumentationURL: "https://cloud.orange-business.com/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/fiches-pratiques/stockage/stockage-objet-s3/guide-de-demarrage/premiere-utilisation-stockage-objet/",
			IPs:              []string{"194.206.55.5"},
			FQDNs:            []string{"s3-region01-priv.cloudavenue.orange-business.com"},
			Ports:            []types.ModelCloudavenueServicesCatalogServicePort{{Port: 443, Protocol: "tcp"}},
		}},
	},
}
