package service

import (
	"context"
	"time"

	"github.com/RIT3shSapata/todo-list-api/internal/couchbase"
	"github.com/RIT3shSapata/todo-list-api/internal/log"
)

type Task struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TaskStatus  string    `json:"taskStatus"`
	UserID      string    `json:"userID"`
	Deadline    time.Time `json:"deadline"`
	TaskID      string    `json:"taskID"`
}

type TaskStatus string

const (
	TaskStatusCreated    TaskStatus = "created"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
)

type CreateTaskPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	UserID      string `json:"userID"`
}

type UpdateTaskPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TaskStatus  string `json:"taskStatus"`
}

type Service struct {
	cluster    couchbase.Cluster
	collection couchbase.Collection
	logger     log.Logger
}

func New(clus couchbase.Cluster, col couchbase.Collection, logger log.Logger) *Service {
	return &Service{
		cluster:    clus,
		collection: col,
		logger:     logger,
	}
}

func (svc *Service) CreateTask(ctx context.Context, payload CreateTaskPayload) (Task, error) {
	return Task{}, nil
}

func (svc *Service) GetTask(ctx context.Context, id string) (Task, error) {
	return Task{}, nil
}

func (svc *Service) UpdateTask(ctx context.Context, id string, payload UpdateTaskPayload) (Task, error) {
	return Task{}, nil
}

func (svc *Service) DeleteTask(ctx context.Context, id string) (Task, error) {
	return Task{}, nil
}
