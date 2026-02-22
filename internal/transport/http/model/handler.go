package model

import "salon/internal/service"

type Handler struct {
	service service.ModelService
}

func NewHandler(s service.ModelService) *Handler {
	return &Handler{
		service: s,
	}
}
