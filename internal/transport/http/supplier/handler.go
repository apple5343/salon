package supplier

import "salon/internal/service"

type Handler struct {
	service service.SupplierService
}

func NewHandler(s service.SupplierService) *Handler {
	return &Handler{
		service: s,
	}
}
