package sale

import (
	"net/http"
	"strconv"
	"time"

	service "salon/internal/models"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

func (h *Handler) GetSales() echo.HandlerFunc {
	return func(c echo.Context) error {
		filters, err := salesFiltersFromRequest(c)
		if err != nil {
			return err
		}
		sales, err := h.service.GetSales(c.Request().Context(), filters)
		if err != nil {
			return err
		}
		result := make([]*models.Sale, len(sales))
		for i, s := range sales {
			result[i] = models.SaleToHttp(s)
		}
		c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(sales)))
		return c.JSON(http.StatusOK, result)
	}
}

func salesFiltersFromRequest(c echo.Context) (*service.SaleFilters, error) {
	var filters service.SaleFilters
	if carID := c.QueryParam("car_id"); carID != "" {
		filters.CarID = &carID
	}
	if status := c.QueryParam("status"); status != "" {
		s, ok := models.SaleStatusMap[status]
		if !ok {
			return nil, errorx.NewError("invalid status", errorx.BadRequest)
		}
		filters.Status = &s
	}
	if paymentType := c.QueryParam("payment_type"); paymentType != "" {
		p, ok := models.PaymentTypeMap[paymentType]
		if !ok {
			return nil, errorx.NewError("invalid payment type", errorx.BadRequest)
		}
		filters.PaymentType = &p
	}
	if finalPriceMin := c.QueryParam("final_price_min"); finalPriceMin != "" {
		p, err := decimal.NewFromString(finalPriceMin)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.FinalPriceMin = &p
	}
	if finalPriceMax := c.QueryParam("final_price_max"); finalPriceMax != "" {
		p, err := decimal.NewFromString(finalPriceMax)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.FinalPriceMax = &p
	}
	if dateFrom := c.QueryParam("date_from"); dateFrom != "" {
		d, err := time.Parse(models.TimeLayout, dateFrom)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.DateFrom = &d
	}
	if dateTo := c.QueryParam("date_to"); dateTo != "" {
		d, err := time.Parse(models.TimeLayout, dateTo)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.DateTo = &d
	}
	if orderBy := c.QueryParam("order_by"); orderBy != "" {
		o, ok := models.SaleOrderByMap[orderBy]
		if !ok {
			return nil, errorx.NewError("invalid order by", errorx.BadRequest)
		}
		filters.OrderBy = &o
	}
	if orderDirection := c.QueryParam("order_direction"); orderDirection != "" {
		o, ok := models.DirectionMap[orderDirection]
		if !ok {
			return nil, errorx.NewError("invalid order direction", errorx.BadRequest)
		}
		filters.OrderDirection = &o
	}
	if limit := c.QueryParam("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.Limit = &l
	}
	if offset := c.QueryParam("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.Offset = &o
	}
	return &filters, nil
}
