package event

import "salon/internal/service"

type Handler struct {
	service service.EventService
}

func NewHandler(s service.EventService) *Handler {
	return &Handler{
		service: s,
	}
}
