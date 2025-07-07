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
	"sync"
)

var mu = &sync.RWMutex{}

type (
	Console      string
	LocationCode string

	console struct {
		SiteName            string
		LocationCode        LocationCode
		SiteID              Console
		Services            Services
		OrganizationPattern *regexp.Regexp
	}

	Services struct {
		IHM         Service
		APIVCD      Service
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
		OrganizationPattern: regexp.MustCompile(`^cav01ev01ocb\d{7}$`),
		Services: Services{
			IHM: Service{
				Enabled:  true,
				Endpoint: "https://console1.cloudavenue.orange-business.com",
			},
			APIVCD: Service{
				Enabled:  true,
				Endpoint: "https://console1.cloudavenue.orange-business.com",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "https://console1.cloudavenue.orange-business.com",
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
		OrganizationPattern: regexp.MustCompile(`^cav01iv02ocb\d{7}$`),
		Services: Services{
			IHM: Service{
				Enabled:  true,
				Endpoint: "https://console2.cloudavenue.orange-business.com",
			},
			APIVCD: Service{
				Enabled:  true,
				Endpoint: "https://console2.cloudavenue.orange-business.com",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "https://console2.cloudavenue.orange-business.com",
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
		OrganizationPattern: regexp.MustCompile(`^cav02ev04ocb\d{7}$`),
		Services: Services{
			IHM: Service{
				Enabled:  true,
				Endpoint: "https://console4.cloudavenue.orange-business.com",
			},
			APIVCD: Service{
				Enabled:  true,
				Endpoint: "https://console4.cloudavenue.orange-business.com",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "https://console4.cloudavenue.orange-business.com",
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
		OrganizationPattern: regexp.MustCompile(`^cav02iv05ocb\d{7}$`),
		Services: Services{
			IHM: Service{
				Enabled:  true,
				Endpoint: "https://console5.cloudavenue-cha.itn.intraorange",
			},
			APIVCD: Service{
				Enabled:  true,
				Endpoint: "https://console5.cloudavenue-cha.itn.intraorange",
			},
			APICerberus: Service{
				Enabled:  true,
				Endpoint: "https://console5.cloudavenue-cha.itn.intraorange",
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
		OrganizationPattern: regexp.MustCompile(`^cav01iv07ocb\d{7}$`),
		Services: Services{
			IHM: Service{
				Enabled:  true,
				Endpoint: "https://console7.cloudavenue-vdr.itn.intraorange",
			},
			APIVCD: Service{
				Enabled:  true,
				Endpoint: "https://console7.cloudavenue-vdr.itn.intraorange",
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
		OrganizationPattern: regexp.MustCompile(`^cav01iv08ocb\d{7}$`),
		Services: Services{
			IHM: Service{
				Enabled:  true,
				Endpoint: "https://console8.cloudavenue-vdr.itn.intraorange",
			},
			APIVCD: Service{
				Enabled:  true,
				Endpoint: "https://console8.cloudavenue-vdr.itn.intraorange",
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
		OrganizationPattern: regexp.MustCompile(`^cav0[0-2]{1}vv09ocb\d{7}$`),
		Services: Services{
			IHM: Service{
				Enabled:  true,
				Endpoint: "https://console9.cloudavenue.orange-business.com",
			},
			APIVCD: Service{
				Enabled:  true,
				Endpoint: "https://console9.cloudavenue.orange-business.com",
			},
			Netbackup: Service{
				Enabled:  false,
				Endpoint: "https://backup9.cloudavenue.orange-business.com/NetBackupSelfService/Api",
			},
		},
	},
}

// FindByOrganizationName - Returns the console by its organization name.
func FindByOrganizationName(organizationName string) (Console, bool) {
	mu.RLock()
	defer mu.RUnlock()

	for c, console := range consoles {
		if console.OrganizationPattern.MatchString(organizationName) {
			return c, true
		}
	}

	return "", false
}

// CheckOrganizationName - Returns true if the organization name is valid.
func CheckOrganizationName(organizationName string) bool {
	if organizationName == "" {
		return false
	}

	mu.RLock()
	defer mu.RUnlock()

	for _, console := range consoles {
		if console.OrganizationPattern.MatchString(organizationName) {
			return true
		}
	}

	return false
}

// Services - Returns the Services.
func (c Console) Services() Services {
	mu.RLock()
	defer mu.RUnlock()

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
	mu.RLock()
	defer mu.RUnlock()

	return consoles[c].SiteName
}

// GetLocationCode - Returns the location code.
func (c Console) GetLocationCode() LocationCode {
	mu.RLock()
	defer mu.RUnlock()

	return consoles[c].LocationCode
}

// GetSiteID - Returns the site ID.
func (c Console) GetSiteID() Console {
	mu.RLock()
	defer mu.RUnlock()

	return consoles[c].SiteID
}

// GetAPIVCDEndpoint - Returns the VMware API endpoint.
func (c Console) GetAPIVCDEndpoint() string {
	mu.RLock()
	defer mu.RUnlock()

	return consoles[c].Services.APIVCD.GetEndpoint()
}

// GetAPICerberusEndpoint - Returns the Cerberus API endpoint.
func (c Console) GetAPICerberusEndpoint() string {
	mu.RLock()
	defer mu.RUnlock()

	return consoles[c].Services.APICerberus.GetEndpoint()
}

// OverrideEndpoint - Overrides the endpoint for a specific service.
func (c Console) OverrideEndpoint(svc Services) {
	mu.Lock()
	defer mu.Unlock()
	x := consoles[c]
	x.Services = svc

	consoles[c] = x
}
