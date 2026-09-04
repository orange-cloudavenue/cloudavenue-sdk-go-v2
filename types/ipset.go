/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

type (
	ParamsListIPSet struct {
		// VDCGroupID is the ID of the Vdc Group owning the IP Sets.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the IP Sets.
		VDCGroupName string
	}

	ParamsGetIPSet struct {
		// ID is the unique identifier of the IP Set.
		ID string

		// Name is the name of the IP Set.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this IP Set.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this IP Set.
		VDCGroupName string
	}

	ParamsCreateIPSet struct {
		// Name is the name of the IP Set.
		Name string

		// Description is the description of the IP Set.
		Description string

		// VDCGroupID is the ID of the Vdc Group owning this IP Set.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this IP Set.
		VDCGroupName string

		// IPAddresses is the list of IP addresses, ranges, or CIDRs (IPv4 or IPv6) in this IP Set.
		IPAddresses []string
	}

	ParamsUpdateIPSet struct {
		// ID is the unique identifier of the IP Set to update.
		ID string

		// Name is the new name of the IP Set.
		Name string

		// Description is the new description of the IP Set.
		Description string

		// VDCGroupID is the ID of the Vdc Group owning this IP Set.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this IP Set.
		VDCGroupName string

		// IPAddresses is the list of IP addresses, ranges, or CIDRs (IPv4 or IPv6) in this IP Set.
		IPAddresses []string
	}

	ParamsDeleteIPSet struct {
		// ID is the unique identifier of the IP Set to delete.
		ID string

		// Name is the name of the IP Set to delete.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this IP Set.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this IP Set.
		VDCGroupName string
	}
)
