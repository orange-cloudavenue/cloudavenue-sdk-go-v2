package commands

import (
	"context"
	"sync"
)

type Registry struct {
	mu       *sync.Mutex
	commands []Command
}

func NewRegistry() *Registry {
	return &Registry{
		mu:       &sync.Mutex{},
		commands: []Command{},
	}
}

func (r *Registry) Register(cmd Command) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands = append(r.commands, cmd)
}

func (r *Registry) Get(namespace, resource, verb string) *Command {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, cmd := range r.commands {
		if cmd.GetNamespace() == namespace && cmd.GetResource() == resource && cmd.GetVerb() == verb {
			return &cmd
		}
	}
	return nil
}

func (c *Command) Run(ctx context.Context, client, params any) (any, error) {
	if len(c.ParamsSpecs) > 0 {
		if err := c.validate(ctx, params); err != nil {
			return nil, err
		}
	}

	v, err := c.RunnerFunc(ctx, c, client, params)
	if err != nil {
		return nil, err
	}
	return v, nil
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
