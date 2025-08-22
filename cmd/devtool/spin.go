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
