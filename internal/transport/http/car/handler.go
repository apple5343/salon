package car

import "salon/internal/service"

type Handler struct {
	service service.CarService
}

func NewHandler(s service.CarService) *Handler {
	return &Handler{
		service: s,
	}
}
