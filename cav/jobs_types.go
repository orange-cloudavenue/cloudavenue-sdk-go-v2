/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

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
		Status JobStatus
	}

	JobStatus string // JobStatus represents the job status, e.g., "queued", "running", "success", "error", "aborted" etc.
)

const (
	JobQueued  JobStatus = "queued"  // Job is queued for execution.
	JobRunning JobStatus = "running" // Job is currently running.
	JobSuccess JobStatus = "success" // Job completed successfully.
	JobError   JobStatus = "error"   // Job encountered an error during execution.
	JobAborted JobStatus = "aborted" // Job was aborted by the user.
)

// IsTerminated checks if the job status is one of the terminal states (Success, Error, Aborted).
func (s JobStatus) IsTerminated() bool {
	return s == JobSuccess || s == JobError || s == JobAborted
}
