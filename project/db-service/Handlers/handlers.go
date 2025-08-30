package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	logger "myproject/project/Logger"
	"myproject/project/db-service/database_connect/service"
	"myproject/project/shared"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	s   service.Service
	log logger.Logger
}

func NewHandler(service service.Service, log logger.Logger) *Handler {
	return &Handler{service, log}
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Invalid ID error(GetTask-handler):%v", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.s.GetTask(ctx, taskID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			h.log.ERROR(fmt.Sprintf("task not found(GetTask-handler): %v", err))
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		h.log.ERROR(fmt.Sprintf("GetTask Handler(db-service) failed: %v", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.log.ERROR(fmt.Sprintf("Json encoding error: %v", err))
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
	h.log.INFO("GetTask Handler executed successfully")
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var ID int
	var erro error

	ctx := r.Context()
	if r.Header.Get("Content-Type") != "application/json" {
		h.log.ERROR(fmt.Sprintf("Wrong Contetnt type in Post Handler(db-service): %v", erro))
		http.Error(w, "Content-Type должен быть application/json", http.StatusUnsupportedMediaType)
		return
	}
	var Task shared.Task
	err := json.NewDecoder(r.Body).Decode(&Task)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Wrong format of JSON in handler(db-service):%v", err))
		http.Error(w, "неверный формат JSON", http.StatusBadRequest)
		return
	}
	h.log.DEBUG(fmt.Sprintf("Post handler: received task: %+v", Task))
	ID, erro = h.s.CreateTask(ctx, Task)
	if erro != nil {
		switch {
		case errors.Is(erro, service.ErrInvalidInput):
			h.log.ERROR(fmt.Sprintf("Post handler: invalid input: %v", err))
			http.Error(w, erro.Error(), http.StatusBadRequest)
		default:
			h.log.ERROR(fmt.Sprintf("Post handler: internal error: %v", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shared.IDResponse{ID: int64(ID)})
	h.log.INFO(fmt.Sprintf("Post handler: task created successfully, ID=%d", ID))
	h.log.DEBUG(fmt.Sprintf("Post handler: response sent with ID=%d", ID))

}
func (h *Handler) AllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.s.GetAllTasks(ctx)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptySlice), errors.Is(err, service.ErrTooFewTasks):

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]shared.Task{})
			return
		default:
			h.log.ERROR(fmt.Sprintf("AllTasks handler:internal error: %v", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
	h.log.INFO("AllTasks handler executed successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	h.log.DEBUG(fmt.Sprintf("Delete handler: got id param = %s", idStr))

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Invalid id: %d", taskID))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.s.ModifyTask(ctx, taskID, "delete")
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			h.log.ERROR(fmt.Sprintf("Delete handler: task %d not found", taskID))
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			h.log.ERROR(fmt.Sprintf("Delete handler: internal error: %v", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	h.log.INFO("Delete handler executed successfully")
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	h.log.DEBUG(fmt.Sprintf("Delete handler: got id param = %s", idStr))

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.ERROR(fmt.Sprintf("Invalid id: %d", taskID))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.s.ModifyTask(ctx, taskID, "updateStatus")
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			h.log.ERROR(fmt.Sprintf("Delete handler: task %d not found", taskID))
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			h.log.ERROR(fmt.Sprintf("Patch handler: internal error: %v", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	h.log.INFO("Patch handler executed successfully")

	w.WriteHeader(http.StatusNoContent)
}
