package service

import (
	"myproject/project/api-service/client"
	"myproject/project/shared"
	//"net/http"
)

type Service struct {
	client *client.Client
}

func NewService(c *client.Client) *Service {
	return &Service{client: c}
}

func (s *Service) Get(id int) (*shared.Task, error) {
	task, err := s.client.GetTask(id)
	if err != nil {
		return nil, err
	}
	return task, nil

}

func (s *Service) GetAll() ([]shared.Task, error) {
	var tasks []shared.Task
	var err error
	tasks, err = s.client.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil

}

func (s *Service) Post(task shared.Task) (int64, error) {
	ID, err := s.client.PostTask(task)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (s *Service) Delete(id int) error {
	err := s.client.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(id int) error {
	err := s.client.Update(id)
	if err != nil {
		return err
	}
	return nil
}
