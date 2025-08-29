package handlers

import (
	"encoding/json"

	"myproject/project/api-service/client"
	"myproject/project/api-service/service"
	"myproject/project/shared"

	//"myproject/project/shared"
	"net/http"
	"strconv"

	"strings"

	"github.com/gorilla/mux"
)

type Handlers struct {
	service service.Service
}

func NewHandler(service service.Service) *Handlers {
	return &Handlers{service}
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Вызов сервиса
	task, err := h.service.Get(taskID)
	if err != nil {
		switch err := err.(type) {
		case *client.NotFoundError:
			http.Error(w, err.Error(), http.StatusNotFound)
		case *client.StatusError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) Post(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "должен быть JSON", http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()

	var task shared.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "неверный формат JSON", http.StatusBadRequest)
		return
	}

	ID, err := h.service.Post(task)
	if err != nil {
		switch e := err.(type) {
		case *client.ContentTypeError:
			http.Error(w, e.Error(), http.StatusUnsupportedMediaType)
		case *client.StatusError:
			http.Error(w, e.Error(), http.StatusConflict)
		default:
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"message": "Задача добавлена",
		"id":      ID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "ошибка кодирования ответа", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAll()
	if err != nil {
		switch e := err.(type) {
		case *client.StatusError:
			http.Error(w, e.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "неверный формат JSON", http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "ошибка кодирования ответа", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(taskID)
	if err != nil {
		switch e := err.(type) {
		case *client.NotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		case *client.StatusError:
			http.Error(w, e.Error(), http.StatusConflict)
		default:
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := shared.DeleteOrUpdateResponse{
		Message:    "Задача успешно удалена",
		StatusCode: http.StatusOK,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "ошибка кодирования ответа", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.Update(taskID)
	if err != nil {
		switch e := err.(type) {
		case *client.NotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		case *client.StatusError:
			http.Error(w, e.Error(), http.StatusConflict)
		default:
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := shared.DeleteOrUpdateResponse{
		Message:    "Задача успешно выполнена ",
		StatusCode: http.StatusOK,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "ошибка кодирования ответа", http.StatusInternalServerError)
		return
	}

}
