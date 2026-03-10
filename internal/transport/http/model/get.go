package model

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
		model, brand, err := h.service.GetByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		if role := ctxutil.UserRoleFromContext(c.Request().Context()); role == string(service.EmployeeRoleAdmin) || role == string(service.EmployeeRoleManager) {
			return c.JSON(http.StatusOK, models.ModelInternalToHttp(model, brand))
		} else {
			return c.JSON(http.StatusOK, models.ModelPublicToHttp(model, brand))
		}
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
		b, ok := models.BodyTypeMap[bodyType]
		if !ok {
			return nil, errorx.NewError("invalid body type", errorx.BadRequest)
		}
		filter.BodyType = &b
	}
	if transmissionType := c.QueryParam("transmission_type"); transmissionType != "" {
		t, ok := models.TransmissionTypeMap[transmissionType]
		if !ok {
			return nil, errorx.NewError("invalid transmission type", errorx.BadRequest)
		}
		filter.TransmissionType = &t
	}
	if fuelType := c.QueryParam("fuel_type"); fuelType != "" {
		f, ok := models.FuelTypeMap[fuelType]
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
	if minPower := c.QueryParam("min_power"); minPower != "" {
		minPowerInt, err := strconv.Atoi(minPower)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MinPowerHP = &minPowerInt
	}
	if maxPower := c.QueryParam("max_power"); maxPower != "" {
		maxPowerInt, err := strconv.Atoi(maxPower)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filter.MaxPowerHP = &maxPowerInt
	}
	if driveType := c.QueryParam("drive_type"); driveType != "" {
		d, ok := models.DriveTypeMap[driveType]
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
	if orderBy := c.QueryParam("order_by"); orderBy != "" {
		o, ok := models.ModelOrderByMap[orderBy]
		if !ok {
			return nil, errorx.NewError("invalid order by", errorx.BadRequest)
		}
		filter.OrderBy = &o
	}

	base, err := models.BaseListFromRequest(c)
	if err != nil {
		return nil, err
	}
	filter.BaseList = base
	return &filter, nil
}
