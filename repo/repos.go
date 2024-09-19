package repo

import (
	"context"
	"poll/models"
)

type RedisRepo interface {
	CreatePoll(ctx context.Context, pollID string, poll models.Poll) error
	GetPoll(ctx context.Context, pollID string) (*models.Poll, error)
	ListPolls(ctx context.Context) ([]models.Poll, error)
	DeletePoll(ctx context.Context, pollID string) error
	UpdatePoll(ctx context.Context, pollID string, poll models.Poll) error
	Close() error
}
