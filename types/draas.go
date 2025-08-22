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
	ModelListDraasOnPremise struct {
		IPs []string `documentation:"List of OnPremise IPs for this organization's draas offer"`
	}

	ParamsAddDraasOnPremiseIP struct {
		IP string
	}

	ParamsRemoveDraasOnPremiseIP struct {
		IP string
	}
)
