package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"poll/configs"
	"poll/models"
	"time"
)

const (
	appID = "poll"
)

var Nil = redis.Nil

type Config struct {
	Timeout time.Duration
}

type RedisRepo struct {
	client *redis.Client
	cfg    configs.RepoConfig
}

func New(ctx context.Context, cfg configs.RepoConfig) (*RedisRepo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		Username: cfg.Redis.Username,
		DB:       cfg.Redis.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("basic connection failure: %v", err)
	}

	return &RedisRepo{
		client: client,
		cfg:    cfg,
	}, nil
}

func (r *RedisRepo) Close() error {
	return r.client.Close()
}

func (s *RedisRepo) generateKey(pollID string) string {
	return fmt.Sprintf("%s:%s", appID, pollID)
}

func (s *RedisRepo) CreatePoll(ctx context.Context, pollID string, poll models.Poll) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, s.cfg.Timeout.Duration)
	defer cancel()

	data, err := json.Marshal(poll)
	if err != nil {
		return fmt.Errorf("failed to marshal poll data: %w", err)
	}

	key := s.generateKey(pollID)
	err = s.client.Set(ctxWithTimeout, key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save poll %s: %w", pollID, err)
	}

	return nil
}

func (s *RedisRepo) GetPoll(ctx context.Context, pollID string) (*models.Poll, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, s.cfg.Timeout.Duration)
	defer cancel()

	key := s.generateKey(pollID)
	data, err := s.client.Get(ctxWithTimeout, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("poll %s not found", pollID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get poll %s: %w", pollID, err)
	}

	var poll models.Poll
	err = json.Unmarshal([]byte(data), &poll)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal poll data: %w", err)
	}

	return &poll, nil
}

func (s *RedisRepo) ListPolls(ctx context.Context) ([]models.Poll, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, s.cfg.Timeout.Duration)
	defer cancel()

	keys, err := s.client.Keys(ctxWithTimeout, "appID:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list polls: %w", err)
	}

	var polls []models.Poll
	for _, key := range keys {
		data, err := s.client.Get(ctxWithTimeout, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get poll for key %s: %w", key, err)
		}

		var poll models.Poll
		err = json.Unmarshal([]byte(data), &poll)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal poll data for key %s: %w", key, err)
		}

		polls = append(polls, poll)
	}

	return polls, nil
}

func (s *RedisRepo) DeletePoll(ctx context.Context, pollID string) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, s.cfg.Timeout.Duration)
	defer cancel()

	key := s.generateKey(pollID)
	err := s.client.Del(ctxWithTimeout, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete poll %s: %w", pollID, err)
	}

	return nil
}

func (s *RedisRepo) UpdatePoll(ctx context.Context, pollID string, poll models.Poll) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, s.cfg.Timeout.Duration)
	defer cancel()

	data, err := json.Marshal(poll)
	if err != nil {
		return fmt.Errorf("failed to marshal poll data: %w", err)
	}

	key := s.generateKey(pollID)
	err = s.client.Set(ctxWithTimeout, key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to update poll %s: %w", pollID, err)
	}

	return nil
}
