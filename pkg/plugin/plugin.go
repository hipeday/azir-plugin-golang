package plugin

import (
	"errors"
	"log"
)

type Plugin interface {
	Run(args []string) (interface{}, error)
}

type CommandFunc func(args []string) (interface{}, error)

type CommandCallbackRenderFunc func(result interface{}, args []string) error

type CommandFunctions struct {
	Command  CommandFunc
	Callback CommandCallbackRenderFunc
}

type CommandRegistry struct {
	commands map[string]CommandFunctions
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]CommandFunctions)}
}

func (r *CommandRegistry) RegisterCommand(name string, fn CommandFunctions) {
	r.commands[name] = fn
}

func (r *CommandRegistry) RunCommand(name string, args []string) error {
	if functions, ok := r.commands[name]; ok {
		result, err := functions.Command(args)
		if err != nil {
			return err
		} else {
			log.Printf("command %s executed successfully", name)
		}
		if functions.Callback != nil {
			err := functions.Callback(result, args)
			if err == nil {
				log.Println("command callback render function executed successfully")
			}
			return err
		} else {
			log.Printf("command %s callback render function is nil so skip it", name)
			return nil
		}
	}
	return errors.New("command not found")
}
