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
	// * List
	ModelListVdcGroup struct {
		VdcGroups []ModelGetVdcGroup `documentation:"List of Vdc Groups"`
	}

	// * Get
	ModelGetVdcGroup struct {
		ID          string `documentation:"ID of the Vdc Group"`
		Name        string `documentation:"Name of the Vdc Group"`
		Description string `documentation:"Description of the Vdc Group"`

		NumberOfVdcs int                   `documentation:"Number of Vdcs in the Vdc Group"`
		Vdcs         []ModelGetVdcGroupVdc `documentation:"List of Vdcs in the Vdc Group"`
	}

	ModelGetVdcGroupVdc struct {
		ID   string `documentation:"ID of the Vdc"`
		Name string `documentation:"Name of the Vdc"`
	}
)

type (
	ParamsListVdcGroup struct {
		// ID is the unique identifier of the Vdc Group to filter by.
		ID string

		// Name is the name of the Vdc Group to filter by.
		Name string
	}

	ParamsGetVdcGroup struct {
		// ID is the unique identifier of the Vdc Group to filter by.
		ID string

		// Name is the name of the Vdc Group to filter by.
		Name string
	}

	ParamsCreateVdcGroup struct {
		// Name is the name of the Vdc Group.
		Name string

		// Description is the description of the Vdc Group.
		Description string

		// Vdcs is the list of Vdcs to associate with the Vdc Group.
		Vdcs []ParamsCreateVdcGroupVdc
	}

	ParamsCreateVdcGroupVdc struct {
		// ID is the unique identifier of the Vdc to associate with the Vdc Group.
		ID string

		// Name is the name of the Vdc to associate with the Vdc Group.
		Name string
	}

	ParamsUpdateVdcGroup struct {
		// ID is the unique identifier of the Vdc Group to update.
		ID string

		// Name is the new name of the Vdc Group.
		Name string

		// Description is the new description of the Vdc Group.
		Description string
	}

	ParamsDeleteVdcGroup struct {
		// ID is the unique identifier of the Vdc Group to delete.
		ID string

		// Name is the name of the Vdc Group to delete.
		Name string

		// Force Value "true" means to forcefully delete the object that contains other objects even if those objects are in a state that does not allow removal. The default is "false"; therefore, objects are not removed if they are not in a state that normally allows removal. Force also implies recursive delete where other contained objects are removed. Errors may be ignored. Invalid value (not true or false) are ignored.
		Force bool
	}

	ParamsAddVdcToVdcGroup struct {
		// ID is the unique identifier of the Vdc Group to update.
		ID string
		// Name is the name of the Vdc Group.
		Name string

		// Vdcs is the list of Vdcs to associate with the Vdc Group.
		Vdcs []ParamsCreateVdcGroupVdc
	}

	ParamsRemoveVdcFromVdcGroup struct {
		// ID is the unique identifier of the Vdc Group to update.
		ID string
		// Name is the name of the Vdc Group.
		Name string

		// Vdcs is the list of Vdcs to disassociate from the Vdc Group.
		Vdcs []ParamsCreateVdcGroupVdc
	}
)
