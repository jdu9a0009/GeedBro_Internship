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
	log := logger.NewLogger("mini-project", logger.LevelInfo)
	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	h := handler.NewHandler(strg, log)

	r := api.NewServer(h)
	r.Run(cfg.Port)
}
