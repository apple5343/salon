package employee

import (
	"net/http"
	service "salon/internal/models"
	"salon/internal/transport/http/models"
	"strconv"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetEmployees() echo.HandlerFunc {
	return func(c echo.Context) error {
		filters, err := employeeFiltersFromRequest(c)
		if err != nil {
			return err
		}
		employeesList, err := h.service.GetEmployees(c.Request().Context(), filters)
		if err != nil {
			return err
		}
		c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(employeesList)))
		result := make([]*models.Employee, len(employeesList))
		for i, e := range employeesList {
			result[i] = models.EmployeeToHttp(e)
		}
		return c.JSON(http.StatusOK, result)
	}
}

func employeeFiltersFromRequest(c echo.Context) (*service.EmployeeFilters, error) {
	var filters service.EmployeeFilters
	if fullName := c.QueryParam("full_name"); fullName != "" {
		filters.FullName = &fullName
	}
	if role := c.QueryParam("role"); role != "" {
		r, err := models.ToServiceRole(role)
		if err != nil {
			return nil, err
		}
		filters.Role = &r
	}
	if status := c.QueryParam("status"); status != "" {
		s, err := models.ToServiceStatus(status)
		if err != nil {
			return nil, err
		}
		filters.Status = &s
	}
	if orderBy := c.QueryParam("order_by"); orderBy != "" {
		o, ok := models.EmployeeOrderByMap[orderBy]
		if !ok {
			return nil, errorx.NewError("invalid order by", errorx.BadRequest)
		}
		filters.OrderBy = &o
	}
	base, err := models.BaseListFromRequest(c)
	if err != nil {
		return nil, err
	}
	filters.BaseList = base

	return &filters, nil
}
