package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"context"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"a1/internal/handlers"
	"a1/internal/queue"
	"a1/internal/repo"
	"a1/internal/services"
)

func main() {
	// 1) Config
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	// 2) DB
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	// 3) Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := pingRedis(rdb); err != nil {
		log.Fatal("cannot connect to redis:", err)
	}

	// 4) Wiring layers
	jobRepo := repo.NewJobRepository(db)
	jobQueue := queue.NewRedisJobQueue(rdb, "jobs_queue")
	jobService := services.NewJobService(jobRepo, jobQueue)
	jobHandler := handlers.NewJobHandler(jobService)

	// 5) Router (versioned API)
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(api chi.Router) {
		jobHandler.RegisterRoutes(api)
	})

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("A1 API listening on", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func pingRedis(rdb *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	return err
}
