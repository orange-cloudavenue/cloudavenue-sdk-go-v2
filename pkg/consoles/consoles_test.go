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
	"testing"
)

func TestConsoles(t *testing.T) {
	tests := []struct {
		name    string
		console Console
		orgName string
		wantErr bool
	}{
		{
			name:    "should return the correct site name for Console1",
			console: Console1,
			orgName: "cav01ev01ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console2",
			console: Console2,
			orgName: "cav01iv02ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console4",
			console: Console4,
			orgName: "cav02ev04ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console5",
			console: Console5,
			orgName: "cav02iv05ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console7",
			console: Console7,
			orgName: "cav01iv07ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console8",
			console: Console8,
			orgName: "cav01iv08ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console9",
			console: Console9,
			orgName: "cav00vv09ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console9",
			console: Console9,
			orgName: "cav01vv09ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console9",
			console: Console9,
			orgName: "cav02vv09ocb0001234",
			wantErr: false,
		},

		{
			name:    "should return an error if the organization is empty",
			console: "",
			orgName: "",
			wantErr: true,
		},
		{
			name:    "should return an error if the organization is invalid",
			console: "",
			orgName: "cav10ev01ocb0001234",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := FindByOrganizationName(tt.orgName)
			if (err != nil) && !tt.wantErr {
				t.Errorf("FingByOrganizationName(%s) error = %v, wantErr %v", tt.orgName, err, tt.wantErr)
				return
			}

			if !tt.wantErr && !CheckOrganizationName(tt.orgName) {
				t.Errorf("CheckOrganizationName(%s) error = %v, wantErr %v", tt.orgName, err, tt.wantErr)
				return
			}

			if c.GetSiteID() != tt.console {
				t.Errorf("FingByOrganizationName(%s) = %v, want %v", tt.orgName, c.GetSiteID(), tt.console)
				return
			}
		})
	}
}

func TestFindBySiteID(t *testing.T) {
	tests := []struct {
		name     string
		siteID   string
		expected Console
		found    bool
	}{
		{
			name:     "valid siteID Console1",
			siteID:   "console1",
			expected: Console1,
			found:    true,
		},
		{
			name:     "valid siteID Console5",
			siteID:   "console5",
			expected: Console5,
			found:    true,
		},
		{
			name:     "invalid siteID",
			siteID:   "invalid",
			expected: "",
			found:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, ok := FindBySiteID(tt.siteID)
			if ok != tt.found {
				t.Errorf("FindBySiteID(%s) found = %v, want %v", tt.siteID, ok, tt.found)
			}
			if c != tt.expected {
				t.Errorf("FindBySiteID(%s) = %v, want %v", tt.siteID, c, tt.expected)
			}
		})
	}
}

func TestFindByURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected Console
		found    bool
	}{
		{
			name:     "valid url Console1",
			url:      "https://console1.cloudavenue.orange-business.com",
			expected: Console1,
			found:    true,
		},
		{
			name:     "valid url Console5",
			url:      "https://console5.cloudavenue-cha.itn.intraorange",
			expected: Console5,
			found:    true,
		},
		{
			name:     "invalid url",
			url:      "https://notfound.example.com",
			expected: "",
			found:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, ok := FindByURL(tt.url)
			if ok != tt.found {
				t.Errorf("FindByURL(%s) found = %v, want %v", tt.url, ok, tt.found)
			}
			if c != tt.expected {
				t.Errorf("FindByURL(%s) = %v, want %v", tt.url, c, tt.expected)
			}
		})
	}
}

func TestConsole_Services(t *testing.T) {
	c := Console1
	services := c.Services()
	if !services.APIVmware.Enabled {
		t.Errorf("Expected APIVmware to be enabled for Console1")
	}
	if services.APIVmware.Endpoint == "" {
		t.Errorf("Expected APIVmware endpoint to be non-empty for Console1")
	}
}

func TestService_IsEnabled(t *testing.T) {
	s := Service{Enabled: true}
	if !s.IsEnabled() {
		t.Errorf("Expected service to be enabled")
	}
	s = Service{Enabled: false}
	if s.IsEnabled() {
		t.Errorf("Expected service to be disabled")
	}
}

func TestService_GetEndpoint(t *testing.T) {
	endpoint := "https://example.com"
	s := Service{Endpoint: endpoint}
	if s.GetEndpoint() != endpoint {
		t.Errorf("Expected endpoint %s, got %s", endpoint, s.GetEndpoint())
	}
}

func TestConsole_GetSiteName(t *testing.T) {
	c := Console4
	expected := "Console Externe CHA"
	if c.GetSiteName() != expected {
		t.Errorf("Expected site name %s, got %s", expected, c.GetSiteName())
	}
}

func TestConsole_GetLocationCode(t *testing.T) {
	c := Console5
	expected := LocationCHR
	if c.GetLocationCode() != expected {
		t.Errorf("Expected location code %s, got %s", expected, c.GetLocationCode())
	}
}

func TestConsole_GetURL(t *testing.T) {
	tests := []struct {
		console  Console
		expected string
	}{
		{Console1, "https://console1.cloudavenue.orange-business.com"},
		{Console5, "https://console5.cloudavenue-cha.itn.intraorange"},
		{Console9, "https://console9.cloudavenue.orange-business.com"},
	}
	for _, tt := range tests {
		got := tt.console.GetURL()
		if got != tt.expected {
			t.Errorf("GetURL() = %v, want %v", got, tt.expected)
		}
	}
}

func TestConsole_GetAPIVmwareEndpoint(t *testing.T) {
	tests := []struct {
		console  Console
		expected string
	}{
		{Console1, "https://console1.cloudavenue.orange-business.com/cloudapi"},
		{Console5, "https://console5.cloudavenue-cha.itn.intraorange/cloudapi"},
		{Console9, "https://console9.cloudavenue.orange-business.com/cloudapi"},
	}
	for _, tt := range tests {
		got := tt.console.GetAPIVmwareEndpoint()
		if got != tt.expected {
			t.Errorf("GetAPIVmwareEndpoint() = %v, want %v", got, tt.expected)
		}
	}
}

func TestConsole_GetAPICerberusEndpoint(t *testing.T) {
	tests := []struct {
		console  Console
		expected string
	}{
		{Console1, "https://console1.cloudavenue.orange-business.com/api/customers"},
		{Console2, "https://console2.cloudavenue.orange-business.com/api/customers"},
		{Console4, "https://console4.cloudavenue.orange-business.com/api/customers"},
		{Console5, "https://console5.cloudavenue-cha.itn.intraorange/api/customers"},
		// Console7, Console8, Console9 do not have Cerberus API enabled, so expect an empty string
		{Console7, ""},
		{Console8, ""},
		{Console9, ""},
	}
	for _, tt := range tests {
		got := tt.console.GetAPICerberusEndpoint()
		if got != tt.expected {
			t.Errorf("GetAPICerberusEndpoint() = %v, want %v", got, tt.expected)
		}
	}
}

func TestCheckOrganizationName(t *testing.T) {
	tests := []struct {
		name string
		org  string
		want bool
	}{
		{
			name: "valid org for Console1",
			org:  "cav01ev01ocb1234567",
			want: true,
		},
		{
			name: "valid org for Console2",
			org:  "cav01iv02ocb7654321",
			want: true,
		},
		{
			name: "valid org for Console4",
			org:  "cav02ev04ocb0000001",
			want: true,
		},
		{
			name: "valid org for Console5",
			org:  "cav02iv05ocb9999999",
			want: true,
		},
		{
			name: "valid org for Console7",
			org:  "cav01iv07ocb1234567",
			want: true,
		},
		{
			name: "valid org for Console8",
			org:  "cav01iv08ocb7654321",
			want: true,
		},
		{
			name: "valid org for Console9 (cav00)",
			org:  "cav00vv09ocb1234567",
			want: true,
		},
		{
			name: "valid org for Console9 (cav01)",
			org:  "cav01vv09ocb7654321",
			want: true,
		},
		{
			name: "valid org for Console9 (cav02)",
			org:  "cav02vv09ocb0000001",
			want: true,
		},
		{
			name: "invalid org (wrong prefix)",
			org:  "cav10ev01ocb1234567",
			want: false,
		},
		{
			name: "invalid org (empty)",
			org:  "",
			want: false,
		},
		{
			name: "invalid org (too short)",
			org:  "cav01ev01ocb1234",
			want: false,
		},
		{
			name: "invalid org (random string)",
			org:  "foobar",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckOrganizationName(tt.org)
			if got != tt.want {
				t.Errorf("CheckOrganizationName(%q) = %v, want %v", tt.org, got, tt.want)
			}
		})
	}
}
