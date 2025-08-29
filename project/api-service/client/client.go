package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myproject/project/shared"
	"net/http"
	"strings"
	"time"
)


type NotFoundError struct {
	Msg string
}

func (e *NotFoundError) Error() string {
	return e.Msg
}

type StatusError struct {
	Code int
	Msg  string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("status %d: %s", e.Code, e.Msg)
}

type ContentTypeError struct {
	Got string
}

func (e *ContentTypeError) Error() string {
	return fmt.Sprintf("unexpected content type: %s", e.Got)
}


type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
	}
}


func (cli *Client) GetTask(id int) (*shared.Task, error) {
	url := fmt.Sprintf("%s/tasks/%d", cli.baseURL, id)

	resp, err := cli.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil, &ContentTypeError{Got: resp.Header.Get("Content-Type")}
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, &NotFoundError{Msg: fmt.Sprintf("task %d not found", id)}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	var task shared.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (cli *Client) PostTask(task shared.Task) (int64, error) {
	var ID shared.IDResponse
	body, err := json.Marshal(task)
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s/tasks", cli.baseURL)
	resp, err := cli.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return 0, &ContentTypeError{Got: resp.Header.Get("Content-Type")}
	}

	if resp.StatusCode == http.StatusConflict {
		return 0, &StatusError{Code: resp.StatusCode, Msg: "task already exists"}
	}

	if resp.StatusCode != http.StatusCreated {
		return 0, &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	if err := json.NewDecoder(resp.Body).Decode(&ID); err != nil {
		return 0, err
	}

	return ID.ID, nil
}

func (cli *Client) GetAllTasks() ([]shared.Task, error) {
	url := fmt.Sprintf("%s/tasks", cli.baseURL)
	resp, err := cli.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	var tasks []shared.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (cli *Client) Delete(id int) error {
	url := fmt.Sprintf("%s/tasks/%d", cli.baseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{Msg: fmt.Sprintf("task %d not found", id)}
	}

	if resp.StatusCode != http.StatusOK {
		return &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	return nil
}

func (cli *Client) Update(id int) error {
	url := fmt.Sprintf("%s/tasks/%d", cli.baseURL, id)
	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		return err
	}

	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{Msg: fmt.Sprintf("task %d not found", id)}
	}

	if resp.StatusCode != http.StatusNoContent {
		return &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	return nil
}
