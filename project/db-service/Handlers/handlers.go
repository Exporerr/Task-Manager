package handlers

import (
	"encoding/json"
	"errors"
	"myproject/project/db-service/database_connect/service"
	"myproject/project/shared"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	s service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	idStr := vars["id"]

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.s.GetTask(ctx, taskID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var ID int
	var erro error

	ctx := r.Context()
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type должен быть application/json", http.StatusUnsupportedMediaType)
		return
	}
	var Task shared.Task
	err := json.NewDecoder(r.Body).Decode(&Task)
	if err != nil {
		http.Error(w, "неверный формат JSON", http.StatusBadRequest)
		return
	}
	ID, erro = h.s.CreateTask(ctx, Task)
	if erro != nil {
		switch {
		case errors.Is(erro, service.ErrInvalidInput):
			http.Error(w, erro.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shared.IDResponse{ID: int64(ID)})

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
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.s.ModifyTask(ctx, taskID, "delete")
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.s.ModifyTask(ctx, taskID, "updateStatus")
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
