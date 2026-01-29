package employee

import "salon/internal/service"

type Handler struct {
	service service.EmployeeService
}

func NewHandler(s service.EmployeeService) *Handler {
	return &Handler{
		service: s,
	}
}
