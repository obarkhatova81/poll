package main

import (
	"context"
	"fmt"
	"log"
	"poll/configs"
	"poll/models"
	"poll/repo/redis"
	httpServer "poll/server/http"
	"poll/server/websocket"
	"poll/service/basic"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	redisClient, err := redis.New(ctx, config.Repo)
	if err != nil {
		log.Fatalf("failed to create Redis client: %v", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalf("failed to close Redis client: %v", err)
		}
	}()

	results := make(chan models.PollResults)
	defer close(results)

	pollService := basic.NewService(redisClient, results)

	httpSrv := httpServer.NewServer(log.Default(), pollService)
	go func() {
		if err := httpSrv.Start(ctx); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	wsSrv := websocket.New(log.Default(), results)
	go func() {
		if err := wsSrv.Start(fmt.Sprintf("0.0.0.0:%s", config.Srv.Monitoring.WebSocket.Port)); err != nil {
			log.Fatalf("WebSocket server failed: %v", err)
		}
	}()

	<-ctx.Done()
}
