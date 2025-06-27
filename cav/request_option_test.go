/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/jobs"
)

// mockJobOption is a dummy JobOption for testing.
func mockJobOption(*jobs.JobOptions) error {
	return nil
}

func TestWithJob_SetsJobOpts(t *testing.T) {
	opt := WithJob(mockJobOption)
	ro := &requestOption{}
	err := opt(ro)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ro.JobOpts == nil {
		t.Fatalf("expected JobOpts to be set, got nil")
	}
}

func TestWithJob_ReturnsErrorOnJobOptionsError(t *testing.T) {
	badJobOpt := func(*jobs.JobOptions) error {
		return assert.AnError // ou toute erreur factice
	}

	opt := WithJob(badJobOpt)
	ro := &requestOption{}
	err := opt(ro)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestWithJob_ReturnsErrorIfJobOptionsFails(t *testing.T) {
	badOpt := func(*requestOption) error {
		return assert.AnError
	}

	_, err := newRequestOptions(badOpt)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
