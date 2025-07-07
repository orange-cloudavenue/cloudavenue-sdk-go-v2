/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package consoles provides utilities for managing and retrieving information
// about Cloud Avenue consoles, their locations, and associated services.
//
// It defines types and constants representing the different consoles and their
// locations, as well as the services available on each console. The package
// includes functions to look up consoles by organization name, validate
// organization names, and retrieve service endpoints and metadata.
//
// Example usage:
//
//	// Find the console for a given organization name
//	console, ok := consoles.FindByOrganizationName("cav01ev01ocb1234567")
//	if ok {
//	    endpoint := console.GetAPIVCDEndpoint()
//	    // Use the endpoint...
//	}
//
// Thread safety:
// All exported functions and methods are safe for concurrent use by multiple goroutines.

package consoles
