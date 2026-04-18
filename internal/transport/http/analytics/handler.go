package analytics

import "salon/internal/service"

type Handler struct {
	service service.AnalyticsService
}

func NewHandler(s service.AnalyticsService) *Handler {
	return &Handler{
		service: s,
	}
}
