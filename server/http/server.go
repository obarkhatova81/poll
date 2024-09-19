package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"poll/service"
)

const (
	listenPort = "8080"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(log *log.Logger, srv service.PollService) *Server {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://your-frontend-domain.com"}, // Замените на ваши домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-Subject-Token", "Link"},
		AllowCredentials: true,
	}))

	h := NewHandler(log, srv)
	h.RegisterRoutes(r)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + listenPort,
			Handler: r,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}
