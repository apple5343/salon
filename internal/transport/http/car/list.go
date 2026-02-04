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

func (h *Handler) GetModels() echo.HandlerFunc {
	return func(c echo.Context) error {
		filters, err := modelFiltersFromRequest(c)
		if err != nil {
			return err
		}
		modelsList, err := h.service.GetModels(c.Request().Context(), filters)
		if err != nil {
			return err
		}
		result := make([]*models.ModelShort, len(modelsList))
		for i, m := range modelsList {
			result[i] = models.ModelShortToHttp(m)
		}
		c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(modelsList)))
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
		s, ok := models.CarStatusType[status]
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
	if offset := c.QueryParam("offset"); offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.Offset = &offsetInt
	}
	if limit := c.QueryParam("limit"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.Limit = &limitInt
	}
	return &filter, nil
}

func modelFiltersFromRequest(c echo.Context) (*service.ModelFilters, error) {
	filter := service.ModelFilters{}
	if brandID := c.QueryParam("brand_id"); brandID != "" {
		filter.BrandID = &brandID
	}
	if name := c.QueryParam("name"); name != "" {
		filter.Name = &name
	}
	if generation := c.QueryParam("generation"); generation != "" {
		filter.Generation = &generation
	}
	if bodyType := c.QueryParam("body_type"); bodyType != "" {
		b, ok := models.BodyType[bodyType]
		if !ok {
			return nil, errorx.NewError("invalid body type", errorx.BadRequest)
		}
		filter.BodyType = &b
	}
	if transmissionType := c.QueryParam("transmission_type"); transmissionType != "" {
		t, ok := models.TransmissionType[transmissionType]
		if !ok {
			return nil, errorx.NewError("invalid transmission type", errorx.BadRequest)
		}
		filter.TransmissionType = &t
	}
	if fuelType := c.QueryParam("fuel_type"); fuelType != "" {
		f, ok := models.FuelType[fuelType]
		if !ok {
			return nil, errorx.NewError("invalid fuel type", errorx.BadRequest)
		}
		filter.FuelType = &f
	}
	if minEngineDisplacement := c.QueryParam("min_engine_displacement"); minEngineDisplacement != "" {
		minEngineDisplacementInt, err := strconv.Atoi(minEngineDisplacement)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MinEngineDisplacement = &minEngineDisplacementInt
	}
	if maxEngineDisplacement := c.QueryParam("max_engine_displacement"); maxEngineDisplacement != "" {
		maxEngineDisplacementInt, err := strconv.Atoi(maxEngineDisplacement)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MaxEngineDisplacement = &maxEngineDisplacementInt
	}
	if driveType := c.QueryParam("drive_type"); driveType != "" {
		d, ok := models.DriveType[driveType]
		if !ok {
			return nil, errorx.NewError("invalid drive type", errorx.BadRequest)
		}
		filter.DriveType = &d
	}
	if minBasePrice := c.QueryParam("min_base_price"); minBasePrice != "" {
		minBasePriceDec, err := strconv.Atoi(minBasePrice)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MinBasePrice = &minBasePriceDec
	}
	if maxBasePrice := c.QueryParam("max_base_price"); maxBasePrice != "" {
		maxBasePriceDec, err := strconv.Atoi(maxBasePrice)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MaxBasePrice = &maxBasePriceDec
	}
	if offset := c.QueryParam("offset"); offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.Offset = &offsetInt
	}
	if limit := c.QueryParam("limit"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.Limit = &limitInt
	}
	if orderBy := c.QueryParam("order_by"); orderBy != "" {
		orderMap := map[string]service.ModelOrderBy{
			"name":                service.ModelOrderByName,
			"base_price":          service.ModelOrderByBasePrice,
			"engine_displacement": service.ModelOrderByEngineDisplacement,
			"power_hp":            service.ModelOrderByPowerHP,
		}
		o, ok := orderMap[orderBy]
		if !ok {
			return nil, errorx.NewError("invalid order by", errorx.BadRequest)
		}
		filter.OrderBy = &o
	}
	if orderDirection := c.QueryParam("order_direction"); orderDirection != "" {
		o, ok := models.DirectionMap[orderDirection]
		if !ok {
			return nil, errorx.NewError("invalid order direction", errorx.BadRequest)
		}
		filter.OrderDirection = &o
	}
	return &filter, nil
}
