package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"poll/models"
)

type Server struct {
	upgrader       websocket.Upgrader
	connections    map[*websocket.Conn]bool
	resultsChannel <-chan models.PollResults
	logger         *log.Logger
}

func New(logger *log.Logger, resultsChannel <-chan models.PollResults) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		connections:    make(map[*websocket.Conn]bool),
		resultsChannel: resultsChannel,
		logger:         logger,
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Printf("failed to upgrade connection: %v", err)
		return
	}

	s.connections[conn] = true

	go func() {
		defer func() {
			if closeErr := conn.Close(); closeErr != nil {
				s.logger.Printf("error closing connection: %v", closeErr)
			}
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				s.logger.Printf("error reading message: %v", err)
				break
			}
		}

		delete(s.connections, conn)
	}()
}

func (s *Server) handleResults() {
	for result := range s.resultsChannel {
		msg, err := json.Marshal(result)
		if err != nil {
			s.logger.Printf("error marshaling poll results: %v", err)
			continue
		}

		for conn := range s.connections {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				s.logger.Printf("error writing message: %v", err)
				conn.Close()
				delete(s.connections, conn)
			}
		}
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/ws", s.handleConnections)
	go s.handleResults()

	s.logger.Printf("websocket server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("websocket server failed: %w", err)
	}

	return nil
}
