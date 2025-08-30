package service

import (
	"fmt"
	logger "myproject/project/Logger"
	"myproject/project/api-service/client"
	"myproject/project/shared"
)

type Service struct {
	client *client.Client
	log    *logger.Logger
}

func NewService(c *client.Client, log *logger.Logger) *Service {
	return &Service{client: c, log: log}
}

func (s *Service) Get(id int) (*shared.Task, error) {
	s.log.DEBUG(fmt.Sprintf("Service: Get task id=%d", id))
	task, err := s.client.GetTask(id)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("Service: Get task failed: %v", err))
		return nil, err
	}
	s.log.INFO(fmt.Sprintf("Service: Get task executed successfully, id=%d", id))
	return task, nil
}

func (s *Service) GetAll() ([]shared.Task, error) {
	s.log.DEBUG("Service: GetAll tasks")
	tasks, err := s.client.GetAllTasks()
	if err != nil {
		s.log.ERROR(fmt.Sprintf("Service: GetAll tasks failed: %v", err))
		return nil, err
	}
	s.log.INFO(fmt.Sprintf("Service: GetAll executed successfully, tasks_count=%d", len(tasks)))
	return tasks, nil
}

func (s *Service) Post(task shared.Task) (int64, error) {
	s.log.DEBUG(fmt.Sprintf("Service: Post task %+v", task))
	ID, err := s.client.PostTask(task)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("Service: Post task failed: %v", err))
		return 0, err
	}
	s.log.INFO(fmt.Sprintf("Service: Post task executed successfully, ID=%d", ID))
	return ID, nil
}

func (s *Service) Delete(id int) error {
	s.log.DEBUG(fmt.Sprintf("Service: Delete task id=%d", id))
	err := s.client.Delete(id)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("Service: Delete task failed: %v", err))
		return err
	}
	s.log.INFO(fmt.Sprintf("Service: Delete task executed successfully, id=%d", id))
	return nil
}

func (s *Service) Update(id int) error {
	s.log.DEBUG(fmt.Sprintf("Service: Update task id=%d", id))
	err := s.client.Update(id)
	if err != nil {
		s.log.ERROR(fmt.Sprintf("Service: Update task failed: %v", err))
		return err
	}
	s.log.INFO(fmt.Sprintf("Service: Update task executed successfully, id=%d", id))
	return nil
}
