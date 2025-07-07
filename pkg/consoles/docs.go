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
