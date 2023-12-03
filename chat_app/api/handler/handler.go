package handler

import (
	"user/pkg/logger"
	"user/storage"
)

type Handler struct {
	storage storage.StorageI
	hub     *Hub
	log     logger.LoggerI
}

func NewHandler(strg storage.StorageI, hub *Hub, loger logger.LoggerI) *Handler {
	return &Handler{storage: strg, hub: hub, log: loger}
}
