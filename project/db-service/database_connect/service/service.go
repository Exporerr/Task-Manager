package service

import (
	"context"
	logger "myproject/project/Logger"
	"myproject/project/db-service/repository"
	"myproject/project/shared"

	"errors"
	"fmt"
	"strings"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrEmptySlice = errors.New("there is no tasks")
var ErrTooFewTasks = errors.New("there are too few tasks")
var ErrInvalidInput = errors.New("invalid input")

type Service struct {
	repo repository.TaskRepository
	log  *logger.Logger
}

func NewService(r repository.TaskRepository, log *logger.Logger) *Service {
	return &Service{r, log}
}

func (s *Service) CreateTask(ctx context.Context, task shared.Task) (int, error) {
	// Валидация входных данных
	if strings.TrimSpace(task.Title) == "" {
		s.log.ERROR(fmt.Sprintf("CreateTask validation failed: title is empty | %v", ErrInvalidInput))
		return 0, fmt.Errorf("%w: title cannot be empty", ErrInvalidInput)
	}
	if task.Status {
		s.log.ERROR(fmt.Sprintf("CreateTask validation failed: status=true | %v", ErrInvalidInput))
		return 0, fmt.Errorf("%w: wrong status: task cannot be created with status = true", ErrInvalidInput)
	}

	// Вызов репозитория
	id, err := s.repo.AddTask(ctx, task)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("CreateTask repo.AddTask failed: %v", err))
		return 0, err
	}

	// Логирование успешного результата
	s.log.INFO(fmt.Sprintf("Task created successfully: ID=%d", id))
	s.log.DEBUG(fmt.Sprintf("CreateTask details: %+v", task))

	return id, nil
}

func (s *Service) GetTask(ctx context.Context, taskID int) (shared.Task, error) {
	//Вызов репозитория 
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("repo.GetTask failed: %v",err))
		return task, err
	}
	//Валидация данных
	if strings.TrimSpace(task.Title) == "" {
		s.log.ERROR(fmt.Sprintf("GetTask validation failed: title is empty|%v",ErrInvalidInput))
		return task, fmt.Errorf("task title is empty")
	}
	s.log.INFO("GetTask(db-service) executed successfully")
	s.log.DEBUG("Success")
	return task, nil
}
func (s *Service) GetAllTasks(ctx context.Context) ([]shared.Task, error) {
	tasks, err := s.repo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, ErrEmptySlice
	}

	if len(tasks) == 1 {
		return nil, ErrTooFewTasks
	}

	return tasks, nil
}

func (s *Service) ModifyTask(ctx context.Context, taskID int, action string) error {
	var rowsAffected int64
	var err error

	switch action {
	case "updateStatus":
		rowsAffected, err = s.repo.UpdateTaskStatus(ctx, taskID)
	case "delete":
		rowsAffected, err = s.repo.DeleteTask(ctx, taskID)
	default:
		return fmt.Errorf("unknown action")
	}

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
