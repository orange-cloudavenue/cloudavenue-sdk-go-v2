package jobs

type (
	// Job struct defines the job status.
	Job struct {
		// ID is the unique identifier of the job.
		ID string

		// Name is the name of the job.
		Name string

		// Description is a brief description of the job.
		Description string

		// HREF is the URL to the job resource.
		HREF string

		// Status is the current status of the job.
		Status Status
	}

	Status string // Status represents the job status, e.g., "queued", "running", "success", "error", "aborted" etc.
)

const (
	Queued  Status = "queued"  // Job is queued for execution.
	Running Status = "running" // Job is currently running.
	Success Status = "success" // Job completed successfully.
	Error   Status = "error"   // Job encountered an error during execution.
	Aborted Status = "aborted" // Job was aborted by the user.
)

// IsTerminated checks if the job status is one of the terminal states (Success, Error, Aborted).
func (s Status) IsTerminated() bool {
	return s == Success || s == Error || s == Aborted
}
