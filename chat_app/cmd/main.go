package main

import (
	"context"
	"fmt"
	"user/api"
	"user/api/handler"
	"user/config"
	"user/pkg/logger"
	"user/storage/postgres"
)

func main() {
	fmt.Println("start")
	cfg := config.Load()
	log := logger.NewLogger("Chat-App", logger.LevelInfo)
	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	hub := handler.NewHub()
	go hub.Run()
	h := handler.NewHandler(strg, hub, log)

	r := api.NewServer(h)
	r.Run(fmt.Sprintf(":%s", cfg.Port))
}
