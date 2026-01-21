package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"yourproject/internal/services"
)

type JobHandler struct {
	Service *services.JobService
}

func NewJobHandler(s *services.JobService) *JobHandler {
	return &JobHandler{Service: s}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req services.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	jobID, err := h.Service.CreateJob(r.Context(), req)
	if err != nil {
		http.Error(w, "failed to submit job", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"job_id": jobID.String(),
		"status": "queued",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	job, err := h.Service.GetJob(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}
