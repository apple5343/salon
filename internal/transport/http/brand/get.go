package brand

import (
	"net/http"
	service "salon/internal/models"
	"salon/internal/transport/http/models"
	ctxutil "salon/internal/utils/context"
	"strconv"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		brand, err := h.service.GetByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		if role := ctxutil.UserRoleFromContext(c.Request().Context()); role == string(service.EmployeeRoleAdmin) || role == string(service.EmployeeRoleManager) {
			return c.JSON(http.StatusOK, models.BrandInternalToHttp(brand))
		}
		return c.JSON(http.StatusOK, models.BrandPublicToHttp(brand))
	}
}

func (h *Handler) GetBrands() echo.HandlerFunc {
	return func(c echo.Context) error {
		filters, err := brandFiltersFromRequest(c)
		if err != nil {
			return err
		}
		brandsList, err := h.service.GetBrands(c.Request().Context(), filters)
		if err != nil {
			return err
		}
		c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(brandsList)))
		if role := ctxutil.UserRoleFromContext(c.Request().Context()); role == string(service.EmployeeRoleAdmin) || role == string(service.EmployeeRoleManager) {
			result := make([]*models.BrandInternalResponse, len(brandsList))
			for i, b := range brandsList {
				result[i] = models.BrandInternalToHttp(b)
			}
			return c.JSON(http.StatusOK, result)
		}
		result := make([]*models.BrandPublicResponse, len(brandsList))
		for i, b := range brandsList {
			result[i] = models.BrandPublicToHttp(b)
		}
		return c.JSON(http.StatusOK, result)
	}
}

func brandFiltersFromRequest(c echo.Context) (*service.BrandFilters, error) {
	var filters service.BrandFilters
	if name := c.QueryParam("name"); name != "" {
		filters.Name = &name
	}
	if countryCode := c.QueryParam("country_code"); countryCode != "" {
		filters.CountryCode = &countryCode
	}
	if orderBy := c.QueryParam("order_by"); orderBy != "" {
		o, ok := models.BrandOrderByMap[orderBy]
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
