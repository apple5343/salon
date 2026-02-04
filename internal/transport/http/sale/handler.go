package sale

import "salon/internal/service"

type Handler struct {
	service service.SaleService
}

func NewHandler(s service.SaleService) *Handler {
	return &Handler{
		service: s,
	}
}
