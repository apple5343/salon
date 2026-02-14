package supplier

import (
	"net/http"
	service "salon/internal/models"
	"salon/internal/transport/http/models"
	ctxutil "salon/internal/utils/context"
	"strconv"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetSuppliers() echo.HandlerFunc {
	return func(c echo.Context) error {
		filters, err := supplierFiltersFromRequest(c)
		if err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		suppliers, err := h.service.GetSuppliers(c.Request().Context(), filters)
		if err != nil {
			return err
		}
		c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(suppliers)))
		if role := ctxutil.UserRoleFromContext(c.Request().Context()); role == string(service.EmployeeRoleAdmin) || role == string(service.EmployeeRoleManager) {
			result := make([]*models.SupplierInternalResponse, len(suppliers))
			for i, s := range suppliers {
				result[i] = models.SupplierInternalToHttp(s)
			}
			return c.JSON(http.StatusOK, result)
		}
		result := make([]*models.SupplierPublicResponse, len(suppliers))
		for i, s := range suppliers {
			result[i] = models.SupplierPublicToHttp(s)
		}
		return c.JSON(http.StatusOK, result)
	}
}

func supplierFiltersFromRequest(c echo.Context) (*service.SupplierFilters, error) {
	var filters service.SupplierFilters
	if name := c.QueryParam("name"); name != "" {
		filters.Name = &name
	}
	if countryCode := c.QueryParam("country_code"); countryCode != "" {
		filters.CountryCode = &countryCode
	}
	if orderBy := c.QueryParam("order_by"); orderBy != "" {
		o, ok := models.SupplierOrderByMap[orderBy]
		if !ok {
			return nil, errorx.NewError("invalid order by", errorx.BadRequest)
		}
		filters.OrderBy = &o
	}
	if orderDirection := c.QueryParam("order_direction"); orderDirection != "" {
		o, ok := models.OrderDirectionMap[orderDirection]
		if !ok {
			return nil, errorx.NewError("invalid order direction", errorx.BadRequest)
		}
		filters.OrderDirection = &o
	}
	if offset := c.QueryParam("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.Offset = &o
	}
	if limit := c.QueryParam("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.Limit = &l
	}
	return &filters, nil
}
