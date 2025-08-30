package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	logger "myproject/project/Logger"
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
	log        *logger.Logger
}

func NewClient(baseURL string, logger logger.Logger) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
		log:        &logger,
	}
}

func (cli *Client) GetTask(id int) (*shared.Task, error) {
	url := fmt.Sprintf("%s/tasks/%d", cli.baseURL, id)
	cli.log.DEBUG(fmt.Sprintf("GET request URL: %s", url)) // DEBUG: формирование запроса

	resp, err := cli.httpClient.Get(url)
	if err != nil {
		cli.log.ERROR(fmt.Sprintf("GET request failed: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		cli.log.ERROR(fmt.Sprintf("unexpected content type: %s", contentType))
		return nil, &ContentTypeError{Got: contentType}
	}

	if resp.StatusCode == http.StatusNotFound {
		cli.log.INFO(fmt.Sprintf("task %d not found", id)) // INFO: ожидаемое отсутствие задачи
		return nil, &NotFoundError{Msg: fmt.Sprintf("task %d not found", id)}
	}

	if resp.StatusCode != http.StatusOK {
		cli.log.ERROR(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
		return nil, &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	var task shared.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		cli.log.ERROR(fmt.Sprintf("failed to decode JSON: %v", err))
		return nil, err
	}

	cli.log.INFO(fmt.Sprintf("task %d retrieved successfully", task.ID)) // INFO: успешное выполнение
	return &task, nil
}

func (cli *Client) PostTask(task shared.Task) (int64, error) {
	body, err := json.Marshal(task)
	if err != nil {
		cli.log.ERROR(fmt.Sprintf("failed to marshal task: %v", err))
		return 0, err
	}

	url := fmt.Sprintf("%s/tasks", cli.baseURL)
	cli.log.DEBUG(fmt.Sprintf("POST request URL: %s, body: %s", url, string(body)))

	resp, err := cli.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		cli.log.ERROR(fmt.Sprintf("POST request failed: %v", err))
		return 0, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		cli.log.ERROR(fmt.Sprintf("unexpected content type: %s", contentType))
		return 0, &ContentTypeError{Got: contentType}
	}

	if resp.StatusCode == http.StatusConflict {
		cli.log.INFO("task already exists") // INFO: ожидаемая конфликтная ситуация
		return 0, &StatusError{Code: resp.StatusCode, Msg: "task already exists"}
	}

	if resp.StatusCode != http.StatusCreated {
		cli.log.ERROR(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
		return 0, &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	var ID shared.IDResponse
	if err := json.NewDecoder(resp.Body).Decode(&ID); err != nil {
		cli.log.ERROR(fmt.Sprintf("failed to decode JSON: %v", err))
		return 0, err
	}

	cli.log.INFO(fmt.Sprintf("task created successfully, ID: %d", ID.ID))
	return ID.ID, nil
}

func (cli *Client) GetAllTasks() ([]shared.Task, error) {
	url := fmt.Sprintf("%s/tasks", cli.baseURL)
	cli.log.DEBUG(fmt.Sprintf("GET ALL request URL: %s", url))

	resp, err := cli.httpClient.Get(url)
	if err != nil {
		cli.log.ERROR(fmt.Sprintf("GET ALL request failed: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		cli.log.ERROR(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
		return nil, &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	var tasks []shared.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		cli.log.ERROR(fmt.Sprintf("failed to decode JSON: %v", err))
		return nil, err
	}

	cli.log.INFO(fmt.Sprintf("retrieved %d tasks successfully", len(tasks)))
	return tasks, nil
}

func (cli *Client) Delete(id int) error {
	url := fmt.Sprintf("%s/tasks/%d", cli.baseURL, id)
	cli.log.DEBUG(fmt.Sprintf("DELETE request URL: %s", url))

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		cli.log.ERROR(fmt.Sprintf("failed to create DELETE request: %v", err))
		return err
	}

	resp, err := cli.httpClient.Do(req)
	if err != nil {
		cli.log.ERROR(fmt.Sprintf("DELETE request failed: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		cli.log.INFO(fmt.Sprintf("task %d not found for deletion", id))
		return &NotFoundError{Msg: fmt.Sprintf("task %d not found", id)}
	}

	if resp.StatusCode != http.StatusOK {
		cli.log.ERROR(fmt.Sprintf("unexpected status code on DELETE: %d", resp.StatusCode))
		return &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	cli.log.INFO(fmt.Sprintf("task %d deleted successfully", id))
	return nil
}

func (cli *Client) Update(id int) error {
	url := fmt.Sprintf("%s/tasks/%d", cli.baseURL, id)
	cli.log.DEBUG(fmt.Sprintf("PATCH request URL: %s", url))

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		cli.log.ERROR(fmt.Sprintf("failed to create PATCH request: %v", err))
		return err
	}

	resp, err := cli.httpClient.Do(req)
	if err != nil {
		cli.log.ERROR(fmt.Sprintf("PATCH request failed: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		cli.log.INFO(fmt.Sprintf("task %d not found for update", id))
		return &NotFoundError{Msg: fmt.Sprintf("task %d not found", id)}
	}

	if resp.StatusCode != http.StatusNoContent {
		cli.log.ERROR(fmt.Sprintf("unexpected status code on PATCH: %d", resp.StatusCode))
		return &StatusError{Code: resp.StatusCode, Msg: "unexpected status"}
	}

	cli.log.INFO(fmt.Sprintf("task %d updated successfully", id))
	return nil
}
