package basic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"poll/models"
	"poll/repo"
)

type PollService struct {
	repo           repo.RedisRepo
	resultsChannel chan<- models.PollResults
}

func NewService(repo repo.RedisRepo, resultsChannel chan<- models.PollResults) *PollService {
	return &PollService{
		repo:           repo,
		resultsChannel: resultsChannel,
	}
}

func (s *PollService) CreatePoll(ctx context.Context, poll models.Poll) (uuid.UUID, error) {
	pollID := uuid.New()

	poll.ID = pollID

	existingPoll, err := s.repo.GetPoll(ctx, pollID.String())
	if err == nil && existingPoll != nil {
		return uuid.Nil, fmt.Errorf("poll with ID %s already exists", pollID.String())
	}

	err = s.repo.CreatePoll(ctx, pollID.String(), poll)
	if err != nil {
		return uuid.Nil, err
	}

	return pollID, nil
}

func (s *PollService) GetPoll(ctx context.Context, pollID string) (*models.Poll, error) {
	poll, err := s.repo.GetPoll(ctx, pollID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving poll: %w", err)
	}
	return poll, nil
}

func (s *PollService) ListPolls(ctx context.Context) ([]models.Poll, error) {
	polls, err := s.repo.ListPolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing polls: %w", err)
	}
	return polls, nil
}

func (s *PollService) DeletePoll(ctx context.Context, pollID string) error {
	_, err := s.repo.GetPoll(ctx, pollID)
	if err != nil {
		return fmt.Errorf("error retrieving poll before deletion: %w", err)
	}

	return s.repo.DeletePoll(ctx, pollID)
}

func (s *PollService) UpdatePoll(ctx context.Context, pollID string, poll models.Poll) error {
	existingPoll, err := s.repo.GetPoll(ctx, pollID)
	if err != nil {
		return fmt.Errorf("error retrieving poll before update: %w", err)
	}
	if existingPoll == nil {
		return fmt.Errorf("poll with ID %s does not exist", pollID)
	}

	if err := s.repo.UpdatePoll(ctx, pollID, poll); err != nil {
		return fmt.Errorf("error updating poll: %w", err)
	}

	return nil
}

func (s *PollService) Vote(ctx context.Context, pollID string, option string) error {
	poll, err := s.repo.GetPoll(ctx, pollID)
	if err != nil {
		return fmt.Errorf("error retrieving poll: %w", err)
	}
	if poll == nil {
		return fmt.Errorf("poll with ID %s does not exist", pollID)
	}

	validOption := false
	for _, o := range poll.Options {
		if o == option {
			validOption = true
			break
		}
	}

	if !validOption {
		return fmt.Errorf("invalid option: %s", option)
	}

	if _, exists := poll.Votes[option]; !exists {
		poll.Votes[option] = 0
	}
	poll.Votes[option]++

	if err := s.repo.UpdatePoll(ctx, pollID, *poll); err != nil {
		return fmt.Errorf("error updating poll: %w", err)
	}

	pollResults := models.PollResults{
		PollID:   pollID,
		Question: poll.Question,
		Options:  poll.Options,
		Votes:    poll.Votes,
	}

	s.resultsChannel <- pollResults

	return nil
}
