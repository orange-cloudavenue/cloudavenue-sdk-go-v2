/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package consoles

import (
	"regexp"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

type (
	Console      string
	LocationCode string

	console struct {
		SiteName            string
		LocationCode        LocationCode
		SiteID              Console
		URL                 string
		Services            Services
		OrganizationPattern *regexp.Regexp
	}

	Services struct {
		APIVmware   Service
		APICerberus Service
		S3          Service
		VCDA        Service
		Netbackup   Service
	}

	Service struct {
		Enabled  bool
		Endpoint string
	}
)

const (
	Console1 Console = "console1" // Externe VDR
	Console2 Console = "console2" // Internal VDR
	Console4 Console = "console4" // Externe CHA
	Console5 Console = "console5" // Internal CHA
	Console7 Console = "console7" // Externe VDR
	Console8 Console = "console8" // Internal VDR
	Console9 Console = "console9" // Externe VDRCHA

	LocationVDR    LocationCode = "vdr"
	LocationCHR    LocationCode = "chr"
	LocationVDRCHA LocationCode = "vdr-cha"
)

var consoles = map[Console]console{
	Console1: {
		SiteName:            "Console Externe VDR",
		LocationCode:        LocationVDR,
		SiteID:              Console1,
		URL:                 "https://console1.cloudavenue.orange-business.com",
		OrganizationPattern: regexp.MustCompile(`^cav01ev01ocb\d{7}$`),
		Services: Services{
			APIVmware: Service{
				Enabled:  true,
				Endpoint: "https://console1.cloudavenue.orange-business.com/cloudapi",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "https://console1.cloudavenue.orange-business.com/api/customers",
			},
			S3: Service{
				Enabled:  true,
				Endpoint: "https://s3console1.cloudavenue.orange-business.com",
			},
			Netbackup: Service{
				Enabled:  true,
				Endpoint: "https://backup1.cloudavenue.orange-business.com/NetBackupSelfService/Api",
			},
		},
	},
	Console2: {
		SiteName:            "Console Interne VDR",
		LocationCode:        LocationVDR,
		SiteID:              Console2,
		URL:                 "https://console2.cloudavenue.orange-business.com",
		OrganizationPattern: regexp.MustCompile(`^cav01iv02ocb\d{7}$`),
		Services: Services{
			APIVmware: Service{
				Enabled:  true,
				Endpoint: "https://console2.cloudavenue.orange-business.com/cloudapi",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "https://console2.cloudavenue.orange-business.com/api/customers",
			},
			S3: Service{
				Enabled:  true,
				Endpoint: "https://s3console2.cloudavenue.orange-business.com",
			},
			Netbackup: Service{
				Enabled:  true,
				Endpoint: "https://backup2.cloudavenue.orange-business.com/NetBackupSelfService/Api",
			},
		},
	},

	Console4: {
		SiteName:            "Console Externe CHA",
		LocationCode:        LocationCHR,
		SiteID:              Console4,
		URL:                 "https://console4.cloudavenue.orange-business.com",
		OrganizationPattern: regexp.MustCompile(`^cav02ev04ocb\d{7}$`),
		Services: Services{
			APIVmware: Service{
				Enabled:  true,
				Endpoint: "https://console4.cloudavenue.orange-business.com/cloudapi",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "https://console4.cloudavenue.orange-business.com/api/customers",
			},
			Netbackup: Service{
				Enabled:  true,
				Endpoint: "https://backup4.cloudavenue.orange-business.com/NetBackupSelfService/Api",
			},
		},
	},
	Console5: {
		SiteName:            "Console Interne CHA",
		LocationCode:        LocationCHR,
		SiteID:              Console5,
		URL:                 "https://console5.cloudavenue-cha.itn.intraorange",
		OrganizationPattern: regexp.MustCompile(`^cav02iv05ocb\d{7}$`),
		Services: Services{
			APIVmware: Service{
				Enabled:  true,
				Endpoint: "https://console5.cloudavenue-cha.itn.intraorange/cloudapi",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "https://console5.cloudavenue-cha.itn.intraorange/api/customers",
			},
			Netbackup: Service{
				Enabled:  true,
				Endpoint: "https://backup5.cloudavenue-cha.itn.intraorange/NetBackupSelfService/Api",
			},
		},
	},

	Console7: {
		SiteName:            "Console specific VDR",
		LocationCode:        LocationVDR,
		SiteID:              Console7,
		URL:                 "https://console7.cloudavenue-vdr.itn.intraorange",
		OrganizationPattern: regexp.MustCompile(`^cav01iv07ocb\d{7}$`),
		Services: Services{
			APIVmware: Service{
				Enabled:  true,
				Endpoint: "https://console7.cloudavenue-vdr.itn.intraorange/cloudapi",
			},
			Netbackup: Service{
				Enabled:  true,
				Endpoint: "https://backup7.cloudavenue-vdr.itn.intraorange/NetBackupSelfService/Api",
			},
		},
	},
	Console8: {
		SiteName:            "Console specific VDR",
		LocationCode:        LocationVDR,
		SiteID:              Console8,
		URL:                 "https://console8.cloudavenue-vdr.itn.intraorange",
		OrganizationPattern: regexp.MustCompile(`^cav01iv08ocb\d{7}$`),
		Services: Services{
			APIVmware: Service{
				Enabled:  true,
				Endpoint: "https://console8.cloudavenue-vdr.itn.intraorange/cloudapi",
			},
			Netbackup: Service{
				Enabled:  true,
				Endpoint: "https://backup8.cloudavenue-vdr.itn.intraorange/NetBackupSelfService/Api",
			},
		},
	},

	Console9: {
		SiteName:            "Console VCOD",
		LocationCode:        LocationVDRCHA,
		SiteID:              Console9,
		URL:                 "https://console9.cloudavenue.orange-business.com",
		OrganizationPattern: regexp.MustCompile(`^cav0[0-2]{1}vv09ocb\d{7}$`),
		Services: Services{
			APIVmware: Service{
				Enabled:  true,
				Endpoint: "https://console9.cloudavenue.orange-business.com/cloudapi",
			},
			Netbackup: Service{
				Enabled:  false,
				Endpoint: "https://backup9.cloudavenue.orange-business.com/NetBackupSelfService/Api",
			},
		},
	},
	// Mock Console for testing purposes
	Console("mock"): {
		SiteName:            "Console Mock",
		LocationCode:        LocationCode("mock"),
		SiteID:              Console("mock"),
		URL:                 "http://mock.api",
		OrganizationPattern: regexp.MustCompile(`^mockorg\d+$`),
		Services: Services{
			APIVmware: Service{
				Enabled:  true,
				Endpoint: "http://mock.api/cloudapi",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "http://mock.api/api/customers",
			},
			S3: Service{
				Enabled:  true,
				Endpoint: "http://mock.api:9000",
			},
			Netbackup: Service{
				Enabled:  true,
				Endpoint: "http://mock.api:8080/NetBackupSelfService/Api",
			},
		},
	},
}

// FindBySiteID - Returns the console by its siteID.
func FindBySiteID(siteID string) (Console, bool) {
	for c, console := range consoles {
		if console.SiteID == Console(siteID) {
			return c, true
		}
	}

	return "", false
}

// FindByURL - Returns the console by its URL.
func FindByURL(url string) (Console, bool) {
	for c, console := range consoles {
		if console.URL == url {
			return c, true
		}
	}

	return "", false
}

// FindByOrganizationName - Returns the console by its organization name.
func FindByOrganizationName(organizationName string) (Console, error) {
	for c, console := range consoles {
		if console.OrganizationPattern.MatchString(organizationName) {
			return c, nil
		}
	}

	return "", errors.ErrOrganizationFormatIsInvalid
}

// CheckOrganizationName - Returns true if the organization name is valid.
func CheckOrganizationName(organizationName string) bool {
	for _, console := range consoles {
		if console.OrganizationPattern.MatchString(organizationName) {
			return true
		}
	}

	return false
}

// Services - Returns the Services.
func (c Console) Services() Services {
	return consoles[c].Services
}

// Enabled - Returns true if the service is enabled.
func (ss Service) IsEnabled() bool {
	return ss.Enabled
}

// GetEndpoint - Returns the endpoint.
func (ss Service) GetEndpoint() string {
	return ss.Endpoint
}

// GetSiteName - Returns the site name.
func (c Console) GetSiteName() string {
	return consoles[c].SiteName
}

// GetLocationCode - Returns the location code.
func (c Console) GetLocationCode() LocationCode {
	return consoles[c].LocationCode
}

// GetSiteID - Returns the site ID.
func (c Console) GetSiteID() Console {
	return consoles[c].SiteID
}

// GetURL - Returns the URL.
func (c Console) GetURL() string {
	return consoles[c].URL
}

// GetAPIVmwareEndpoint - Returns the VMware API endpoint.
func (c Console) GetAPIVmwareEndpoint() string {
	return consoles[c].Services.APIVmware.GetEndpoint()
}

// GetAPICerberusEndpoint - Returns the Cerberus API endpoint.
func (c Console) GetAPICerberusEndpoint() string {
	return consoles[c].Services.APICerberus.GetEndpoint()
}
