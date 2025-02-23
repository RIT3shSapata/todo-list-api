package main

import (
	"context"
	"fmt"
	"os"

	"github.com/RIT3shSapata/todo-list-api/cmd/api/server"
	"github.com/RIT3shSapata/todo-list-api/internal/config"
	"github.com/RIT3shSapata/todo-list-api/internal/log"
	"go.uber.org/zap"
)

func main() {

	config, err := config.NewAPIConfig()
	if err != nil {
		fmt.Printf("failed to load config: %s", err)
		os.Exit(1)
	}

	logger, err := log.New(&config.LogOpts)
	if err != nil {
		fmt.Printf("failed to create logger")
		os.Exit(1)
	}

	server, err := server.NewAPI(logger, config)
	if err != nil {
		logger.Error("failed to create api: %w", zap.Error(err))
	}

	err = server.Start(context.Background(), config.ServerPort)
	if err != nil {
		logger.Error("failed to start api: %w", zap.Error(err))
	}

}
