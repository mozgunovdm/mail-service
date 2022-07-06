package db

import (
	"context"
	"mts/task-service/internal/models"
)

type IDatabase interface {
	GetTask(id string) (*models.Task, error)
	WriteTask(models.Task) error
	DeleteTask(id, login string) error
	UpdateTask(models.Task) error
	ListTask() (models.Tasks, error)
	Close(context.Context)
}
