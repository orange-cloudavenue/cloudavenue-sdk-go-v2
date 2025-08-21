/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

import (
	"slices"
	"sync"
)

var globalRegistry = newRegistry()

type Registry struct {
	mu       *sync.RWMutex
	Commands []Command
}

func NewRegistry() *Registry {
	return globalRegistry
}

func newRegistry() *Registry {
	return &Registry{
		mu:       &sync.RWMutex{},
		Commands: []Command{},
	}
}

func (r *Registry) Register(cmd Command) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Commands = append(r.Commands, cmd)
}

func (r *Registry) Get(namespace, resource, verb string) *Command {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, cmd := range r.Commands {
		if cmd.GetNamespace() == namespace && cmd.GetResource() == resource && cmd.GetVerb() == verb {
			return &cmd
		}
	}
	return nil
}

func (r *Registry) GetNamespaces() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	namespaces := make([]string, 0)
	for _, cmd := range r.Commands {
		if !slices.Contains(namespaces, cmd.GetNamespace()) {
			namespaces = append(namespaces, cmd.GetNamespace())
		}
	}
	return namespaces
}

func (r *Registry) GetCommandsByFilter(filter func(cmd Command) bool) []Command {
	r.mu.RLock()
	defer r.mu.RUnlock()
	commands := make([]Command, 0)
	for _, cmd := range r.Commands {
		if filter(cmd) {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func (c *Command) GetNamespace() string {
	return c.Namespace
}

func (c *Command) GetResource() string {
	return c.Resource
}

func (c *Command) GetVerb() string {
	return c.Verb
}
