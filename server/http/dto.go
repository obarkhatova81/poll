package server

import (
	"github.com/google/uuid"
)

type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type UpdatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type PollResponse struct {
	ID       uuid.UUID      `json:"id"`
	Question string         `json:"question"`
	Options  []string       `json:"options"`
	Votes    map[string]int `json:"votes"`
}

type VoteRequest struct {
	Option string `json:"option"`
	UserID string `json:"user_id"`
}
