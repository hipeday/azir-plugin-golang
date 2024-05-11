package plugin

import (
	"encoding/json"
	"flag"
	"github.com/hipeday/azir-plugin-golang/pkg/properties"
	"go.uber.org/zap"
	"log"
)

type Config interface {
	Plugin
	ParseConfig(args []string) (interface{}, error)
	GetConfig() interface{}
	GetLogger() *zap.SugaredLogger
}

type ConfigPlugin struct {
	Config
	conf   interface{}
	logger *zap.SugaredLogger
}

func (c *ConfigPlugin) ParseConfig(args []string) (interface{}, error) {
	var (
		property properties.DefaultProperty
		err      error
	)

	runCommand := flag.NewFlagSet("run", flag.ExitOnError)
	configArg := runCommand.String("c", "", "Configuration JSON")
	err = runCommand.Parse(args)
	if err != nil {
		log.Fatalf("Error parsing args: %v", err)
	}

	err = json.Unmarshal([]byte(*configArg), &property)
	if err != nil {
		return nil, err
	}
	c.conf = property
	return property, err
}

func (c *ConfigPlugin) GetConfig() interface{} {
	return c.conf
}

func (c *ConfigPlugin) GetLogger() *zap.SugaredLogger {
	if c.logger != nil {
		return c.logger
	}
	property := c.GetConfig().(properties.DefaultProperty)
	c.logger = property.Logger.CreateLogger(property.Home, property.InvokeId)
	return c.logger
}
