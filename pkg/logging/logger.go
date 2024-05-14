package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	// Encoding can be one "json" or "console". Defaults to "console"
	Encoding string `json:"encoding"`

	// Level configures the log level
	Level string `json:"level"`

	// Colors configures if color output should be enabled
	Colors *bool `json:"colors"`

	// time format
	TimeFormat string `json:"time_format"`
}

func (c *Config) CreateLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	timeFormat := "2006-01-02 15:04:05.000"

	if c.TimeFormat != "" {
		timeFormat = c.TimeFormat
	}

	// 设置时间格式
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(timeFormat)

	if c != nil {
		if c.Encoding == "" || c.Encoding == "console" {
			config.Encoding = "console"
		}

		if c.Colors != nil && *c.Colors {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		switch c.Level {
		case "debug":
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		case "info":
			config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		case "warn":
			config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "error":
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		}
	}

	logger, _ := config.Build()
	sugar := logger.Sugar()
	return sugar
}
