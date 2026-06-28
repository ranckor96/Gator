package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	cMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cMap[cmd.Name]
	if !ok {
		return fmt.Errorf("%s command not found", cmd.Name)
	}
	return handler(s, cmd)
}
