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
	ParamsListSecurityGroup struct {
		// VdcGroupID is the ID of the Vdc Group owning the Security Groups.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning the Security Groups.
		VdcGroupName string
	}

	ParamsGetSecurityGroup struct {
		// ID is the unique identifier of the Security Group.
		ID string

		// Name is the name of the Security Group.
		Name string

		// VdcGroupID is the ID of the Vdc Group owning this Security Group.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Security Group.
		VdcGroupName string
	}

	ParamsCreateSecurityGroup struct {
		// Name is the name of the Security Group.
		Name string

		// Description is the description of the Security Group.
		Description string

		// VdcGroupID is the ID of the Vdc Group owning this Security Group.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Security Group.
		VdcGroupName string

		// Members is the list of Org VDC network references to attach to this Security Group.
		Members []ParamsFirewallGroupMember
	}

	ParamsFirewallGroupMember struct {
		// ID is the unique identifier of the network.
		ID string

		// Name is the name of the network.
		Name string
	}

	ParamsUpdateSecurityGroup struct {
		// ID is the unique identifier of the Security Group to update.
		ID string

		// Name is the new name of the Security Group.
		Name string

		// Description is the new description of the Security Group.
		Description string

		// VdcGroupID is the ID of the Vdc Group owning this Security Group.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Security Group.
		VdcGroupName string

		// Members is the list of Org VDC network references to attach to this Security Group.
		Members []ParamsFirewallGroupMember
	}

	ParamsDeleteSecurityGroup struct {
		// ID is the unique identifier of the Security Group to delete.
		ID string

		// Name is the name of the Security Group to delete.
		Name string

		// VdcGroupID is the ID of the Vdc Group owning this Security Group.
		VdcGroupID string

		// VdcGroupName is the name of the Vdc Group owning this Security Group.
		VdcGroupName string
	}
)
