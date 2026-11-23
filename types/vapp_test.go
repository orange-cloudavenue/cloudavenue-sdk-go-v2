/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

import "testing"

func TestModelVAppGuestProperties_ValidateRejectsKeyInSeveralSections(t *testing.T) {
	err := ModelVAppGuestProperties{Sections: []ModelVAppGuestPropertySection{
		{Class: "a", Properties: map[string]ModelVAppGuestProperty{"k": {}}},
		{Class: "b", Properties: map[string]ModelVAppGuestProperty{"k": {}}},
	}}.Validate()
	if err == nil {
		t.Fatal("expected an error for a key defined in two sections")
	}
}
