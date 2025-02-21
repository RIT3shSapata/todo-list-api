package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RIT3shSapata/todo-list-api/internal/config"
	"github.com/RIT3shSapata/todo-list-api/internal/couchbase"
	"github.com/RIT3shSapata/todo-list-api/internal/endpoints"
	"github.com/RIT3shSapata/todo-list-api/internal/log"
	"github.com/RIT3shSapata/todo-list-api/internal/tasks"
	tasksSvc "github.com/RIT3shSapata/todo-list-api/internal/tasks/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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

func (api *API) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// g, ctx := errgroup.WithContext(ctx)

	router := mux.NewRouter()

	router.Use(api.loggingMiddleware)

	return nil
}

func (api *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.String()
		method := r.Method

		api.logger.Info("recived request", zap.String("method", method), zap.String("uri", uri))
		next.ServeHTTP(w, r)
	})
}
