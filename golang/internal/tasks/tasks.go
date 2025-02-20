package tasks

import (
	"context"

	tasksSvc "github.com/RIT3shSapata/todo-list-api/internal/tasks/service"
)

type Svc interface {
	GetTask(ctx context.Context, id string) (tasksSvc.Task, error)
	UpdateTask(ctx context.Context, id string, payload tasksSvc.UpdateTaskPayload) (tasksSvc.Task, error)
	CreateTask(ctx context.Context, payload tasksSvc.CreateTaskPayload) (tasksSvc.Task, error)
	DeleteTask(ctx context.Context, id string) (tasksSvc.Task, error)
}
