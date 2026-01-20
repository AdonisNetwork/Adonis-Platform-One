package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"context"
	"time"
)

type Server struct {
	db    *sql.DB
	redis *redis.Client
}

type createJobRequest struct {
	Task  string `json:"task"`
	Input string `json:"input"`
}

func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	s := &Server{db: db, redis: rdb}

	http.HandleFunc("/api/jobs", s.handleCreateJob)
	http.HandleFunc("/api/jobs/", s.handleGetJob)

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Println("API listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Task == "" || req.Input == "" {
		http.Error(w, "task and input required", http.StatusBadRequest)
		return
	}

	id := uuid.New()

	_, err := s.db.Exec(
		`INSERT INTO jobs (id, task_type, input_text, status) VALUES ($1, $2, $3, $4)`,
		id, req.Task, req.Input, "queued",
	)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// push to Redis queue
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.redis.LPush(ctx, "jobs_queue", id.String()).Err(); err != nil {
		http.Error(w, "queue error", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"job_id": id.String(),
		"status": "queued",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/jobs/"):]
	if idStr == "" {
		http.Error(w, "job id required", http.StatusBadRequest)
		return
	}

	var (
		taskType, status, input, resultMarkdown, errorText sql.NullString
		resultJSON                                         []byte
		createdAt, updatedAt                               time.Time
	)
	err := s.db.QueryRow(`
		SELECT task_type, status, input_text, result_markdown, result_json, error_text, created_at, updated_at
		FROM jobs WHERE id = $1
	`, idStr).Scan(&taskType, &status, &input, &resultMarkdown, &resultJSON, &errorText, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	var jsonResult interface{}
	if len(resultJSON) > 0 {
		_ = json.Unmarshal(resultJSON, &jsonResult)
	}

	resp := map[string]interface{}{
		"job_id":          idStr,
		"task_type":       taskType.String,
		"status":          status.String,
		"input":           input.String,
		"result_markdown": resultMarkdown.String,
		"result_json":     jsonResult,
		"error":           errorText.String,
		"created_at":      createdAt,
		"updated_at":      updatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
