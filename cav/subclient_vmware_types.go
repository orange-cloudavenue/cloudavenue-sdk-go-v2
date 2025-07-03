package cav

type (

	// vmware struct is the VMware subclient that implements the Client interface.
	// It provides methods to interact with VMware Cloud Director API.
	vmware struct {
		subclient
	}

	// VmwareError represents the error structure returned by VMware Cloud Director API.
	// It contains a message and a minor error code.
	// This structure is used to parse error responses from the VMware API.
	vmwareError struct {
		// DOC API : https://developer.broadcom.com/xapis/vmware-cloud-director-api/38.1/doc/types/ErrorType.html
		Message       string `json:"message"`
		StatusCode    int    `json:"majorErrorCode"`
		StatusMessage string `json:"minorErrorCode"`
	}
)
