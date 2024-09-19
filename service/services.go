package service

import (
	"context"
	"github.com/google/uuid"
	"poll/models"
)

type PollService interface {
	CreatePoll(ctx context.Context, poll models.Poll) (uuid.UUID, error)
	GetPoll(ctx context.Context, pollID string) (*models.Poll, error)
	ListPolls(ctx context.Context) ([]models.Poll, error)
	DeletePoll(ctx context.Context, pollID string) error
	UpdatePoll(ctx context.Context, pollID string, poll models.Poll) error
	Vote(ctx context.Context, pollID string, option string) error
}
