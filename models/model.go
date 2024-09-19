package models

import "github.com/google/uuid"

type Poll struct {
	ID       uuid.UUID      `json:"id"`
	Question string         `json:"question"`
	Options  []string       `json:"options"`
	Votes    map[string]int `json:"votes"`
}

type PollResults struct {
	PollID   string         `json:"poll_id"`
	Question string         `json:"question"`
	Options  []string       `json:"options"`
	Votes    map[string]int `json:"votes"`
}
