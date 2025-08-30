package handlers

import (
	"encoding/json"
	"fmt"
	logger "myproject/project/Logger"
	"myproject/project/api-service/client"
	"myproject/project/api-service/service"
	"myproject/project/shared"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Handlers struct {
	service service.Service
	log     *logger.Logger
}

func NewHandler(service service.Service, log *logger.Logger) *Handlers {
	return &Handlers{service, log}
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Get handler: invalid id: %v", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	h.log.DEBUG(fmt.Sprintf("Get handler: received id=%d", taskID))

	task, err := h.service.Get(taskID)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Get handler: service error: %v", err))
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

	h.log.INFO(fmt.Sprintf("Get handler: task retrieved successfully, id=%d", taskID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) Post(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		h.log.ERROR("Post handler: wrong content type")
		http.Error(w, "должен быть JSON", http.StatusUnsupportedMediaType)
		return
	}
	defer r.Body.Close()

	var task shared.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.log.ERROR(fmt.Sprintf("Post handler: wrong JSON format: %v", err))
		http.Error(w, "неверный формат JSON", http.StatusBadRequest)
		return
	}

	h.log.DEBUG(fmt.Sprintf("Post handler: received task %+v", task))
	ID, err := h.service.Post(task)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Post handler: service error: %v", err))
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

	h.log.INFO(fmt.Sprintf("Post handler: task created successfully, ID=%d", ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Задача добавлена",
		"id":      ID,
	})
}

func (h *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	h.log.DEBUG("GetAll handler: called")
	tasks, err := h.service.GetAll()
	if err != nil {
		h.log.ERROR(fmt.Sprintf("GetAll handler: service error: %v", err))
		switch e := err.(type) {
		case *client.StatusError:
			http.Error(w, e.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "неверный формат JSON", http.StatusBadRequest)
		}
		return
	}
	h.log.INFO(fmt.Sprintf("GetAll handler executed successfully, tasks_count=%d", len(tasks)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Delete handler: invalid id: %v", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	h.log.DEBUG(fmt.Sprintf("Delete handler: received id=%d", taskID))

	err = h.service.Delete(taskID)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Delete handler: service error: %v", err))
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

	h.log.INFO(fmt.Sprintf("Delete handler executed successfully, id=%d", taskID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shared.DeleteOrUpdateResponse{
		Message:    "Задача успешно удалена",
		StatusCode: http.StatusOK,
	})
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Update handler: invalid id: %v", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	h.log.DEBUG(fmt.Sprintf("Update handler: received id=%d", taskID))

	err = h.service.Update(taskID)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Update handler: service error: %v", err))
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

	h.log.INFO(fmt.Sprintf("Update handler executed successfully, id=%d", taskID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shared.DeleteOrUpdateResponse{
		Message:    "Задача успешно выполнена",
		StatusCode: http.StatusOK,
	})
}
