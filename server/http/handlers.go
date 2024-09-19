package server

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	_ "poll/docs"
	"poll/models"
	"poll/service"
)

type Handler struct {
	log *log.Logger
	srv service.PollService
}

func NewHandler(log *log.Logger, srv service.PollService) *Handler {
	return &Handler{
		log: log,
		srv: srv,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/polls", h.CreatePoll)
	r.Get("/polls/{id}", h.GetPoll)
	r.Put("/polls/{id}", h.UpdatePoll)
	r.Delete("/polls/{id}", h.DeletePoll)
	r.Get("/polls", h.ListPolls)
	r.Post("/polls/{id}/vote", h.VoteHandler)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}

// @Tags Polls
// @Summary Create a new poll
// @Description Create a new poll with a unique ID
// @Produce json
// @Param poll body CreatePollRequest true "Poll data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /polls [post]
func (h *Handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req CreatePollRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	poll := models.Poll{
		Question: req.Question,
		Options:  req.Options,
		Votes:    make(map[string]int),
	}

	pollID, err := h.srv.CreatePoll(r.Context(), poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "Poll created successfully",
		"pollID": pollID,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// @Tags Polls
// @Summary Get a poll by ID
// @Description Retrieve a poll by its unique ID
// @Produce json
// @Param id path string true "Poll ID"
// @Success 200 {object} models.Poll
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /polls/{id} [get]
func (h *Handler) GetPoll(w http.ResponseWriter, r *http.Request) {
	pollID := chi.URLParam(r, "id")

	poll, err := h.srv.GetPoll(r.Context(), pollID)
	if err != nil {
		if err.Error() == "poll not found" {
			http.Error(w, "Poll not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(poll); err != nil {
		http.Error(w, "Failed to encode poll response", http.StatusInternalServerError)
	}
}

// @Tags Polls
// @Summary List all polls
// @Description Retrieve a list of all polls
// @Produce json
// @Success 200 {array} models.Poll
// @Failure 500 {object} map[string]string
// @Router /polls [get]
func (h *Handler) ListPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := h.srv.ListPolls(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(polls); err != nil {
		http.Error(w, "Failed to encode polls response", http.StatusInternalServerError)
		return
	}
}

// @Tags Polls
// @Summary Delete a poll by ID
// @Description Delete a poll by its unique ID
// @Produce json
// @Param id path string true "Poll ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /polls/{id} [delete]
func (h *Handler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	pollID := chi.URLParam(r, "id")

	err := h.srv.DeletePoll(r.Context(), pollID)
	if err != nil {
		if err.Error() == "poll not found" {
			http.Error(w, "Poll not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"status": "Poll deleted successfully",
	}

	w.WriteHeader(http.StatusOK)

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		h.log.Printf("error encoding response: %v", encodeErr)
	}
}

// @Tags Polls
// @Summary Update a poll by ID
// @Description Update a poll's details by its unique ID
// @Produce json
// @Param id path string true "Poll ID"
// @Param poll body UpdatePollRequest true "Poll data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /polls/{id} [put]
func (h *Handler) UpdatePoll(w http.ResponseWriter, r *http.Request) {
	var req UpdatePollRequest
	pollID := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	poll := models.Poll{
		ID:       uuid.MustParse(pollID),
		Question: req.Question,
		Options:  req.Options,
		Votes:    make(map[string]int),
	}

	err := h.srv.UpdatePoll(r.Context(), pollID, poll)
	if err != nil {
		if err.Error() == "poll not found" {
			http.Error(w, "Poll not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"status": "Poll updated successfully",
	}); encodeErr != nil {
		h.log.Printf("error encoding response: %v", encodeErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// VoteHandler handles voting for a poll.
// @Summary Vote for a poll
// @Description Allows a user to vote for a poll option
// @Tags Poll
// @Accept json
// @Produce json
// @Param id path string true "Poll ID"
// @Param vote body VoteRequest true "Vote details"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 404 {object} map[string]string "Poll not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /polls/{id}/vote [post]
func (h *Handler) VoteHandler(w http.ResponseWriter, r *http.Request) {
	pollID := chi.URLParam(r, "id")

	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.srv.Vote(r.Context(), pollID, req.Option)
	if err != nil {
		if err.Error() == "poll not found" {
			http.Error(w, "Poll not found", http.StatusNotFound)
		} else if err.Error() == "option not found" {
			http.Error(w, "Option not found", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"status": "Vote recorded successfully",
	}); encodeErr != nil {
		h.log.Printf("error encoding response: %v", encodeErr)
	}
}
