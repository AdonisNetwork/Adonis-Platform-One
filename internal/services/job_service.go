package services

import (
	"context"

	"github.com/google/uuid"

	"a1/internal/queue"
	"a1/internal/repo"
)

type CreateJobRequest struct {
	Task       string                 `json:"task"`
	Input      string                 `json:"input"`
	InputMeta  map[string]interface{} `json:"input_meta,omitempty"`
}

type JobService struct {
	Repo  *repo.JobRepository
	Queue queue.JobQueue
}

func NewJobService(r *repo.JobRepository, q queue.JobQueue) *JobService {
	return &JobService{
		Repo:  r,
		Queue: q,
	}
}

func (s *JobService) CreateJob(ctx context.Context, req CreateJobRequest) (uuid.UUID, error) {
	// می‌توانیم در آینده validation پیشرفته اضافه کنیم
	id, err := s.Repo.InsertJob(ctx, req.Task, req.Input, req.InputMeta)
	if err != nil {
		return uuid.Nil, err
	}

	// پرتاب به صف worker
	if err := s.Queue.Enqueue(ctx, id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *JobService) GetJob(ctx context.Context, id uuid.UUID) (*repo.Job, error) {
	return s.Repo.GetJob(ctx, id)
}
