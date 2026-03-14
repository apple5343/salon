package car

import (
	"net/http"
	service "salon/internal/models"
	"salon/internal/transport/http/models"
	"strconv"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

func (h *Handler) GetCars() echo.HandlerFunc {
	return func(c echo.Context) error {
		filters, err := carFiltersFromRequest(c)
		if err != nil {
			return err
		}
		cars, err := h.service.GetCars(c.Request().Context(), filters)
		if err != nil {
			return err
		}
		result := make([]*models.CarShort, len(cars))
		for i, c := range cars {
			result[i] = models.CarShortToHttp(c)
		}
		c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(cars)))
		return c.JSON(http.StatusOK, result)
	}
}

func carFiltersFromRequest(c echo.Context) (*service.CarFilters, error) {
	filter := service.CarFilters{}
	if supplierID := c.QueryParam("supplier_id"); supplierID != "" {
		filter.SupplierID = &supplierID
	}
	if modelID := c.QueryParam("model_id"); modelID != "" {
		filter.ModelID = &modelID
	}
	if brandID := c.QueryParam("brand_id"); brandID != "" {
		filter.BrandID = &brandID
	}
	if minYear := c.QueryParam("min_year"); minYear != "" {
		minYearInt, err := strconv.Atoi(minYear)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MinYear = &minYearInt
	}
	if maxYear := c.QueryParam("max_year"); maxYear != "" {
		maxYearInt, err := strconv.Atoi(maxYear)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MaxYear = &maxYearInt
	}
	if color := c.QueryParam("color"); color != "" {
		filter.Color = &color
	}
	if status := c.QueryParam("status"); status != "" {
		s, ok := models.CarStatusTypeMap[status]
		if !ok {
			return nil, errorx.NewError("invalid status", errorx.BadRequest)
		}
		filter.Status = &s
	}
	if minPrice := c.QueryParam("min_price"); minPrice != "" {
		minPriceDec, err := decimal.NewFromString(minPrice)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MinPrice = &minPriceDec
	}
	if maxPrice := c.QueryParam("max_price"); maxPrice != "" {
		maxPriceDec, err := decimal.NewFromString(maxPrice)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MaxPrice = &maxPriceDec
	}
	if minMileage := c.QueryParam("min_mileage"); minMileage != "" {
		minMileageInt, err := strconv.Atoi(minMileage)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MinMileage = &minMileageInt
	}
	if maxMileage := c.QueryParam("max_mileage"); maxMileage != "" {
		maxMileageInt, err := strconv.Atoi(maxMileage)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MaxMileage = &maxMileageInt
	}
	base, err := models.BaseListFromRequest(c)
	if err != nil {
		return nil, err
	}
	filter.BaseList = base
	return &filter, nil
}
