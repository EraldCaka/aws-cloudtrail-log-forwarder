package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

type RedisService struct {
	client      *redis.Client
	ctx         context.Context
	clientAsynq *asynq.Client
	worker      *asynq.Server
}

func NewRedisService() *RedisService {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	clientAsynq := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})

	worker := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	//	mux.HandleFunc("send_email", handleSendEmailTask)

	go func() {
		if err := worker.Run(mux); err != nil {
			log.Fatalf("Could not start worker: %v", err)
		}
	}()

	return &RedisService{
		client:      client,
		ctx:         ctx,
		clientAsynq: clientAsynq,
		worker:      worker,
	}
}

func (r *RedisService) SetKey(key string, value string) error {
	err := r.client.Set(r.ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %v", err)
	}
	return nil
}

func (r *RedisService) GetKey(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key from Redis: %v", err)
	}
	return val, nil
}

func (r *RedisService) DeleteKey(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key from Redis: %v", err)
	}
	return nil
}

func (r *RedisService) EnqueueTask(queue string, to string) error {
	payload := map[string]interface{}{
		"to": to,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	task := asynq.NewTask(queue, payloadBytes)

	_, err = r.clientAsynq.Enqueue(task, asynq.Queue(queue))
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %v", err)
	}
	return nil
}

// func handleSendEmailTask(c context.Context, t *asynq.Task) error {
// 	fmt.Printf("Sending email to: %s\n")
// 	return nil
// }
