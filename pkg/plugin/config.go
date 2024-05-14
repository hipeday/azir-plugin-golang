package plugin

import (
	"encoding/json"
	"flag"
	"github.com/hipeday/azir-plugin-golang/pkg/logging"
	"go.uber.org/zap"
	"log"
)

type CommandParam struct {
	Name    string `json:"name,omitempty"`
	Default string `json:"default,omitempty"`
	Usage   string `json:"usage,omitempty"`
	Value   string
}

type Config interface {
	Plugin
	ParseConfig(command string, params []CommandParam, args []string) error
	GetLogger(args []string) *zap.SugaredLogger
}

type ConfigPlugin struct {
	Config
	logger *zap.SugaredLogger
}

func (c *ConfigPlugin) ParseConfig(command string, params []CommandParam, args []string) error {
	runCommand := flag.NewFlagSet(command, flag.ExitOnError)
	for _, param := range params {
		runCommand.StringVar(&param.Value, param.Name, param.Default, param.Usage)
	}

	return runCommand.Parse(args)
}

func (c *ConfigPlugin) GetLogger(args []string) *zap.SugaredLogger {
	var (
		loggerConf   string
		loggerConfig logging.Config
		err          error
	)

	command := flag.NewFlagSet("", flag.ExitOnError)
	command.StringVar(&loggerConf, "l", "", "logger config")
	err = command.Parse(args)
	if err != nil {
		var colors = true
		// 初始化
		loggerConfig = logging.Config{
			Encoding:   "console",
			Level:      "info",
			Colors:     &colors,
			TimeFormat: "2006-01-02 15:04:05.000",
		}
	} else {
		err = json.Unmarshal([]byte(loggerConf), &loggerConfig)
		if err != nil {
			log.Fatalf("Error parsing logger config: %v", err)
		}
	}
	c.logger = loggerConfig.CreateLogger()
	return c.logger
}
