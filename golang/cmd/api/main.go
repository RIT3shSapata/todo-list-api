package main

import (
	"fmt"
	"os"

	"github.com/RIT3shSapata/todo-list-api/internal/log"
)

func main() {

	logOpts := log.LogOpts{
		Name:    "api",
		Level:   log.Info,
		Encoder: log.LogJSONEncoder,
	}

	logger, err := log.New(&logOpts)
	if err != nil {
		fmt.Printf("failed to create logger")
		os.Exit(1)
	}

	logger.Info("API is working")
	logger.Debug("API is working")
	logger.Error("API is working")
}
