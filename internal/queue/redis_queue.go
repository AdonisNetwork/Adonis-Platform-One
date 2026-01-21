package queue

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JobQueue interface {
	Enqueue(ctx context.Context, jobID uuid.UUID) error
}

type RedisJobQueue struct {
	Client    *redis.Client
	QueueName string
	Timeout   time.Duration
}

func NewRedisJobQueue(client *redis.Client, queueName string) *RedisJobQueue {
	return &RedisJobQueue{
		Client:    client,
		QueueName: queueName,
		Timeout:   2 * time.Second,
	}
}

func (q *RedisJobQueue) Enqueue(ctx context.Context, jobID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, q.Timeout)
	defer cancel()

	return q.Client.LPush(ctx, q.QueueName, jobID.String()).Err()
}
