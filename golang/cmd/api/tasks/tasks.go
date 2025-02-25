package tasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RIT3shSapata/todo-list-api/internal/endpoints"
	"github.com/RIT3shSapata/todo-list-api/internal/log"
	tasksSvc "github.com/RIT3shSapata/todo-list-api/internal/tasks"
	"github.com/gorilla/mux"
)

const (
	root     = "/tasks"
	taskRoot = "/tasks/{taskID}"
)

type Handler struct {
	logger    log.Logger
	responder endpoints.Responder
	tasksSvc  tasksSvc.Svc
}

func New(logger log.Logger, responder endpoints.Responder, tasksSvc tasksSvc.Svc) *Handler {
	return &Handler{
		logger:    logger,
		responder: responder,
		tasksSvc:  tasksSvc,
	}
}

func (h *Handler) GetTask() *endpoints.Endpoint {
	return &endpoints.Endpoint{
		Handler: func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

			taskID := mux.Vars(r)["taskID"]

			task, err := h.tasksSvc.GetTask(ctx, taskID)
			if err != nil {
				h.responder.Error(ctx, w, endpoints.WithError(
					fmt.Errorf("error fetching task: %w", err)),
					endpoints.WithStatusCode(http.StatusInternalServerError),
					endpoints.WithLog())
			}
			h.responder.Respond(ctx, w, task)
		},
	}
}

func (h *Handler) Register(router *mux.Router) {
	router.Handle(taskRoot, h.GetTask()).Methods(http.MethodGet)
}
