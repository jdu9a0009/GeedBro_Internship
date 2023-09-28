package handler

import (
	"user/pkg/logger"
	"user/storage"
)

type Handler struct {
	storage storage.StorageI

	log logger.LoggerI
}

func NewHandler(strg storage.StorageI, loger logger.LoggerI) *Handler {
	return &Handler{storage: strg, log: loger}
}
