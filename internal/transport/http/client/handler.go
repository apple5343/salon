package client

import "salon/internal/service"

type Handler struct {
	service service.ClientService
}

func NewHandler(s service.ClientService) *Handler {
	return &Handler{
		service: s,
	}
}
