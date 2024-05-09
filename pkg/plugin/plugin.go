package plugin

import (
	"errors"
)

type Plugin interface {
	Run(args []string) (interface{}, error)
}

type CommandFunc func(args []string) (interface{}, error)

type CommandRegistry struct {
	commands map[string]CommandFunc
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]CommandFunc)}
}

func (r *CommandRegistry) RegisterCommand(name string, fn CommandFunc) {
	r.commands[name] = fn
}

func (r *CommandRegistry) RunCommand(name string, args []string) (interface{}, error) {
	if function, ok := r.commands[name]; ok {
		return function(args)
	}
	return "", errors.New("command not found")
}
