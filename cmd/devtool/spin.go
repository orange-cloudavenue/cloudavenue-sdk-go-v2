/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package main

import (
	"fmt"
	"time"
)

var (
	monkeys = []string{"🙈", "🙈", "🙉", "🙊"}
)

func spinner(message string, animation []string, interval time.Duration) func() {
	ch := make(chan struct{})
	i := 0
	go func() {
		for {
			select {
			case <-time.After(interval):
				fmt.Printf("\r%s %s", animation[i], message)
				i = (i + 1) % len(animation)
			case <-ch:
				fmt.Printf("\r%s %s\n", "✔", message)
				return
			}
		}
	}()
	return func() { close(ch) }
}
