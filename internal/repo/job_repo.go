package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID             uuid.UUID              `json:"job_id"`
	TaskType       string                 `json:"task_type"`
	InputText      string                 `json:"input"`
	InputMeta      map[string]interface{} `json:"input_meta,omitempty"`
	Status         string                 `json:"status"`
	StatusReason   string                 `json:"status_reason,omitempty"`
	ResultMarkdown string                 `json:"result_markdown,omitempty"`
	ResultJSON     map[string]interface{} `json:"result_json,omitempty"`
	ErrorText      string                 `json:"error,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type JobRepository struct {
	DB *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{DB: db}
}

func (r *JobRepository) InsertJob(ctx context.Context, taskType, inputText string, inputMeta map[string]interface{}) (uuid.UUID, error) {
	id := uuid.New()

	var metaJSON []byte
	if inputMeta != nil {
		var err error
		metaJSON, err = jsonMarshal(inputMeta)
		if err != nil {
			return uuid.Nil, err
		}
	}

	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO jobs (id, task_type, input_text, input_meta, status)
		VALUES ($1, $2, $3, $4, 'queued')
	`, id, taskType, inputText, metaJSON)
	if err != nil {
		return uuid.Nil, err
	}

	_, _ = r.DB.ExecContext(ctx, `
		INSERT INTO job_events (job_id, event_type, payload)
		VALUES ($1, $2, $3)
	`, id, "queued", metaJSON)

	return id, nil
}

func (r *JobRepository) GetJob(ctx context.Context, id uuid.UUID) (*Job, error) {
	var (
		taskType, status, inputText, statusReason sql.NullString
		resultMarkdown, errorText                sql.NullString
		inputMetaBytes, resultJSONBytes          []byte
		createdAt, updatedAt                     time.Time
	)

	err := r.DB.QueryRowContext(ctx, `
		SELECT task_type,
		       status,
		       input_text,
		       input_meta,
		       status_reason,
		       result_markdown,
		       result_json,
		       error_text,
		       created_at,
		       updated_at
		FROM jobs
		WHERE id = $1
	`, id).Scan(
		&taskType,
		&status,
		&inputText,
		&inputMetaBytes,
		&statusReason,
		&resultMarkdown,
		&resultJSONBytes,
		&errorText,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var meta map[string]interface{}
	if len(inputMetaBytes) > 0 {
		_ = jsonUnmarshal(inputMetaBytes, &meta)
	}

	var resultJSON map[string]interface{}
	if len(resultJSONBytes) > 0 {
		_ = jsonUnmarshal(resultJSONBytes, &resultJSON)
	}

	return &Job{
		ID:             id,
		TaskType:       taskType.String,
		InputText:      inputText.String,
		InputMeta:      meta,
		Status:         status.String,
		StatusReason:   statusReason.String,
		ResultMarkdown: resultMarkdown.String,
		ResultJSON:     resultJSON,
		ErrorText:      errorText.String,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
