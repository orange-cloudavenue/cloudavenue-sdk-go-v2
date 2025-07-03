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
			c, ok := FindByOrganizationName(tt.orgName)
			if !ok && !tt.wantErr {
				t.Errorf("FindByOrganizationName(%s) = %v, want error", tt.orgName, c)
				return
			}

			if !tt.wantErr {
				got := CheckOrganizationName(tt.orgName)
				if !got {
					t.Errorf("CheckOrganizationName(%s) = %v, want %v", tt.orgName, got, !tt.wantErr)
					return
				}
			}

			if c.GetSiteID() != tt.console {
				t.Errorf("FindByOrganizationName(%s) = %v, want %v", tt.orgName, c.GetSiteID(), tt.console)
				return
			}
		})
	}
}

func TestConsole_Services(t *testing.T) {
	c := Console1
	services := c.Services()
	if !services.APIVCD.Enabled {
		t.Errorf("Expected APIVCD to be enabled for Console1")
	}
	if services.APIVCD.Endpoint == "" {
		t.Errorf("Expected APIVCD endpoint to be non-empty for Console1")
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

func TestConsole_GetAPIVCDEndpoint(t *testing.T) {
	tests := []struct {
		console  Console
		expected string
	}{
		{Console1, "https://console1.cloudavenue.orange-business.com/cloudapi"},
		{Console5, "https://console5.cloudavenue-cha.itn.intraorange/cloudapi"},
		{Console9, "https://console9.cloudavenue.orange-business.com/cloudapi"},
	}
	for _, tt := range tests {
		got := tt.console.GetAPIVCDEndpoint()
		if got != tt.expected {
			t.Errorf("GetAPIVCDEndpoint() = %v, want %v", got, tt.expected)
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

func TestConsole_OverrideEndpoint(t *testing.T) {
	original := Console1.Services()
	newServices := original
	newServices.APIVCD.Endpoint = "https://custom-endpoint.example.com/cloudapi"
	newServices.APIVCD.Enabled = false

	Console1.OverrideEndpoint(newServices)

	got := Console1.Services()
	if got.APIVCD.Endpoint != "https://custom-endpoint.example.com/cloudapi" {
		t.Errorf("OverrideEndpoint did not update APIVCD endpoint, got %s", got.APIVCD.Endpoint)
	}
	if got.APIVCD.Enabled != false {
		t.Errorf("OverrideEndpoint did not update APIVCD enabled flag, got %v", got.APIVCD.Enabled)
	}

	// Restore original state for other tests
	Console1.OverrideEndpoint(original)
}
