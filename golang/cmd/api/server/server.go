package server

import (
	"fmt"

	"github.com/RIT3shSapata/todo-list-api/internal/config"
	"github.com/RIT3shSapata/todo-list-api/internal/couchbase"
	"github.com/RIT3shSapata/todo-list-api/internal/endpoints"
	"github.com/RIT3shSapata/todo-list-api/internal/log"
	"github.com/RIT3shSapata/todo-list-api/internal/tasks"
	tasksSvc "github.com/RIT3shSapata/todo-list-api/internal/tasks/service"
)

type API struct {
	logger    log.Logger
	responder endpoints.Responder
	tasksSvc  tasks.Svc
}

func NewAPI(logger log.Logger, cfg config.Config) (*API, error) {
	responder := endpoints.Responder{
		Logger: logger,
	}

	clus, err := couchbase.NewCluster(cfg.CouchbaseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to couchbase: %w", err)
	}

	col := clus.BucketDefaultCol(cfg.CouchbaseConfig.Bucket)

	tasksSvc := tasksSvc.New(clus, col, logger)

	return &API{
		logger:    logger,
		responder: responder,
		tasksSvc:  tasksSvc,
	}, nil
}
