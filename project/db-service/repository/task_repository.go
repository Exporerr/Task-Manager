package repository

import (
	"context"
	databaseconnect "myproject/project/db-service/database_connect"
	"myproject/project/shared"
)

type TaskRepository interface {
	GetTask(ctx context.Context, id int) (shared.Task, error)        //
	AddTask(ctx context.Context, task shared.Task) (int, error)      //
	GetAllTasks(ctx context.Context) ([]shared.Task, error)          //
	UpdateTaskStatus(ctx context.Context, taskID int) (int64, error) //
	DeleteTask(ctx context.Context, taskID int) (int64, error)       //
}

func NewTaskRepository(s *databaseconnect.Storage) TaskRepository {
	return s
}
