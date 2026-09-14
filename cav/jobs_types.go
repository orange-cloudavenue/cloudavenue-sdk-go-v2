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
	// Job describes an asynchronous backend job.
	Job struct {
		// ID is job identifier.
		ID string

		// Name is job name.
		Name string

		// Description summarizes job purpose or outcome.
		Description string

		// HREF is job resource URL when backend exposes one.
		HREF string

		// Status is current job state.
		Status JobStatus
	}

	// JobStatus represents backend job state.
	JobStatus string
)

const (
	// JobQueued indicates job is waiting to run.
	JobQueued JobStatus = "queued"
	// JobRunning indicates job is in progress.
	JobRunning JobStatus = "running"
	// JobSuccess indicates job completed successfully.
	JobSuccess JobStatus = "success"
	// JobError indicates job completed with failure.
	JobError JobStatus = "error"
	// JobAborted indicates job was canceled.
	JobAborted JobStatus = "aborted"
)

// IsTerminated reports whether s is terminal.
func (s JobStatus) IsTerminated() bool {
	return s == JobSuccess || s == JobError || s == JobAborted
}

// String returns s as string.
func (s JobStatus) String() string {
	return string(s)
}
