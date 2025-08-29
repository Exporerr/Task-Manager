package databaseconnect

import (
	"context"
	"fmt"
	logger "myproject/project/Logger"
	"myproject/project/shared"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db  *pgxpool.Pool
	log *logger.Logger
}

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func NewUserPool(pool *pgxpool.Pool, log *logger.Logger) *Storage {
	return &Storage{db: pool, log: log}
}

func (s *Storage) AddTask(ctx context.Context, task shared.Task) (int, error) {
	var insertedID int

	query := `
        INSERT INTO tasks (title, description, status)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	err := s.db.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Status,
	).Scan(&insertedID)

	if err != nil {
		s.log.ERROR(fmt.Sprintf("failed to execute query AddTask: %v", err))
		return 0, err
	}
	s.log.DEBUG(fmt.Sprintf("AddTask executed successfully, ID: %d", insertedID))

	return insertedID, nil
}

func (s *Storage) GetTask(ctx context.Context, id int) (shared.Task, error) {

	var Task shared.Task

	query := `SELECT id, title, description, status, created_at FROM tasks WHERE id = $1`

	err := s.db.QueryRow(ctx, query, id).Scan(
		&Task.ID,
		&Task.Title,
		&Task.Description,
		&Task.Status,
		&Task.Created_at,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			s.log.ERROR(fmt.Sprintf("task with id %d not found", id))
			return Task, fmt.Errorf("task with id %d not found", id)
		}
		s.log.ERROR(fmt.Sprintf("GetTask failed: %v", err))
		return Task, err
	}
	s.log.DEBUG("GetTask executed successfully!")

	return Task, nil
}

func (s *Storage) GetAllTasks(ctx context.Context) ([]shared.Task, error) {
	query := `SELECT id, title, description, status, created_at FROM tasks ORDER BY created_at DESC`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("GetAllTasks failed: %v", err))
		return nil, err
	}
	defer rows.Close()

	tasks := []shared.Task{}

	for rows.Next() {
		var t shared.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Created_at); err != nil {
			s.log.ERROR(fmt.Sprintf("GetAllTasks scan failed:%v", err))
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		s.log.ERROR(fmt.Sprintf("GetAllTasks rows error: %v", err))
		return nil, err
	}
	s.log.INFO(fmt.Sprintf("GetAllTasks executed successfully, count=%d", len(tasks)))
	s.log.DEBUG("GetAllTasks query executed")

	return tasks, nil
}

func (s *Storage) UpdateTaskStatus(ctx context.Context, taskID int) (int64, error) {
	query := `UPDATE tasks SET status=true WHERE id=$1`
	cmdTag, err := s.db.Exec(ctx, query, taskID)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("UpdateTaskStatus failed for ID=%d: %v", taskID, err))
		return 0, err
	}
	s.log.INFO(fmt.Sprintf("Task status updated successfully: ID=%d", taskID))
	s.log.DEBUG(fmt.Sprintf("UpdateTaskStatus query executed for ID=%d", taskID))
	return cmdTag.RowsAffected(), nil
}

func (s *Storage) DeleteTask(ctx context.Context, taskID int) (int64, error) {
	query := `DELETE FROM tasks WHERE id=$1`
	cmdTag, err := s.db.Exec(ctx, query, taskID)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("DeleteTask failed for ID=%d: %v", taskID, err))
		return 0, err
	}
	s.log.INFO(fmt.Sprintf("Task deleted successfully: ID=%d", taskID))
	s.log.DEBUG(fmt.Sprintf("DeleteTask query executed for ID=%d", taskID))
	return cmdTag.RowsAffected(), nil
}
