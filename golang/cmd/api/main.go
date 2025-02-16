package main

import (
	"fmt"
	"os"

	"github.com/RIT3shSapata/todo-list-api/internal/config"
	"github.com/RIT3shSapata/todo-list-api/internal/log"
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

	logger.Info("API is working")
	logger.Debug("API is working")
	logger.Error("API is working")
}
