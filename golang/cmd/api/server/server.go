package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RIT3shSapata/todo-list-api/internal/config"
	"github.com/RIT3shSapata/todo-list-api/internal/couchbase"
	"github.com/RIT3shSapata/todo-list-api/internal/endpoints"
	"github.com/RIT3shSapata/todo-list-api/internal/log"
	"github.com/RIT3shSapata/todo-list-api/internal/tasks"
	tasksSvc "github.com/RIT3shSapata/todo-list-api/internal/tasks/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	Timeout = 2 * time.Second
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

func (api *API) Start(ctx context.Context, port string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	router := mux.NewRouter()

	router.Use(api.loggingMiddleware)

	router = api.registerHandlers(router)

	port = ":" + port

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	api.logger.Info("api is listening on", zap.String("address:", server.Addr))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		sg := <-sig
		api.logger.Info("received signal, shutting down", zap.Any("signal", sg))
		cancel()
	}()

	g.Go(func() error {
		<-ctx.Done()
		// ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		// defer cancel()
		return server.Shutdown(ctx)
	})

	var err error
	g.Go(func() error {
		err = server.ListenAndServe()
		if err != http.ErrServerClosed {
			return fmt.Errorf("api error: %w", err)
		}
		return nil
	})

	return g.Wait()
}

func (api *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.String()
		method := r.Method

		api.logger.Info("recived request", zap.String("method", method), zap.String("uri", uri))
		next.ServeHTTP(w, r)
	})
}

func (api *API) registerHandlers(router *mux.Router) *mux.Router {
	router.HandleFunc("/", api.root)
	return router
}

func (api *API) root(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	if _, err := io.WriteString(w, "todo-api is running"); err != nil {
		api.logger.Error("error getting root", zap.Error(err))
	}
}
