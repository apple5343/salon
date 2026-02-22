package brand

import "salon/internal/service"

type Handler struct {
	service service.BrandService
}

func NewHandler(s service.BrandService) *Handler {
	return &Handler{
		service: s,
	}
}
