package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"a1/internal/services"
)

type JobHandler struct {
	Service *services.JobService
}

func NewJobHandler(s *services.JobService) *JobHandler {
	return &JobHandler{Service: s}
}

func (h *JobHandler) RegisterRoutes(r chi.Router) {
	r.Post("/jobs", h.CreateJob)
	r.Get("/jobs/{id}", h.GetJob)
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req services.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Task == "" || req.Input == "" {
		http.Error(w, "task and input required", http.StatusBadRequest)
		return
	}

	jobID, err := h.Service.CreateJob(r.Context(), req)
	if err != nil {
		http.Error(w, "failed to create job", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"job_id": jobID.String(),
		"status": "queued",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "job id required", http.StatusBadRequest)
		return
	}

	jobID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	job, err := h.Service.GetJob(r.Context(), jobID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(job)
}
