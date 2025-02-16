package config

import (
	"fmt"

	"github.com/RIT3shSapata/todo-list-api/internal/couchbase"
	"github.com/RIT3shSapata/todo-list-api/internal/log"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Database        string `env:"DATABASE"`
	Env             string `env:"GOLANG_API_ENV"`
	LogEncoder      string `env:"GOLANG_API_LOG_ENCODER"`
	CouchbaseConfig couchbase.CouchbaseConfig
	LogOpts         log.LogOpts
}

func NewAPIConfig() (Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	logLevel := log.Debug
	switch config.Env {
	case "dev":
		logLevel = log.Debug
	case "prod":
		logLevel = log.Info
	}

	logEncoder := log.LogJSONEncoder
	switch config.LogEncoder {
	case "json":
		logEncoder = log.LogJSONEncoder
	case "console":
		logEncoder = log.LogConsoleEncoder
	}

	logOpts := log.LogOpts{
		Name:    "api",
		Level:   logLevel,
		Encoder: logEncoder,
	}

	config.LogOpts = logOpts
	return config, nil
}
