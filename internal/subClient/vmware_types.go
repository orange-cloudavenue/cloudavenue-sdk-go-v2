package subclient

// vmware struct is the VMware subclient that implements the Client interface.
// It provides methods to interact with VMware Cloud Director API.
type vmware struct {
	client
}

// vmwareError represents the error structure returned by VMware Cloud Director API.
// It contains a message and a minor error code.
// This structure is used to parse error responses from the VMware API.
type vmwareError struct {
	Message        string `json:"message"`
	MinorErrorCode string `json:"minorErrorCode"`
}

// vmwareJobAPIResponse represents an asynchronous operation in VMware Cloud Director.
type vmwareJobAPIResponse struct {
	HREF             string       `json:"HREF,omitempty"`             // The URI of the entity.
	ID               string       `json:"ID,omitempty"`               // The entity identifier, expressed in URN format. The value of this attribute uniquely identifies the entity, persists for the life of the entity, and is never reused.
	OperationKey     string       `json:"operationKey,omitempty"`     // Optional unique identifier to support idempotent semantics for create and delete operations.
	Name             string       `json:"name,omitempty"`             // The name of the entity.
	Status           string       `json:"status,omitempty"`           // The execution status of the task. One of queued, preRunning, running, success, error, aborted
	Operation        string       `json:"operation,omitempty"`        // A message describing the operation that is tracked by this task.
	OperationName    string       `json:"operationName,omitempty"`    // The short name of the operation that is tracked by this task.
	ServiceNamespace string       `json:"serviceNamespace,omitempty"` // Identifier of the service that created the task. It must not start with com.vmware.vcloud and the length must be between 1 and 128 symbols.
	StartTime        string       `json:"startTime,omitempty"`        // The date and time the system started executing the task. May not be present if the task has not been executed yet.
	EndTime          string       `json:"endTime,omitempty"`          // The date and time that processing of the task was completed. May not be present if the task is still being executed.
	ExpiryTime       string       `json:"expiryTime,omitempty"`       // The date and time at which the task resource will be destroyed and no longer available for retrieval. May not be present if the task has not been executed or is still being executed.
	CancelRequested  bool         `json:"cancelRequested,omitempty"`  // Whether user has requested this processing to be canceled.
	Description      string       `json:"description,omitempty"`      // Optional description.
	Error            *vmwareError `json:"error,omitempty"`            // Represents error information from a failed task.
	Progress         int          `json:"progress,omitempty"`         // Read-only indicator of task progress as an approximate percentage between 0 and 100. Not available for all tasks.
	Details          string       `json:"details,omitempty"`          // Detailed message about the task. Also contained by the Owner entity when task status is preRunning.
}
