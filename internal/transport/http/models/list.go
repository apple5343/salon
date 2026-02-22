package models

import (
	"salon/internal/models"
	"strconv"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

var OrderDirectionMap = map[string]models.OrderDirection{
	"asc":  models.OrderDirectionASC,
	"desc": models.OrderDirectionDESC,
}

func BaseListFromRequest(c echo.Context) (models.BaseList, error) {
	var filters models.BaseList
	if offset := c.QueryParam("offset"); offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return filters, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.Offset = &offsetInt
	}
	if limit := c.QueryParam("limit"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return filters, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.Limit = &limitInt
	}
	if orderDirection := c.QueryParam("order_direction"); orderDirection != "" {
		o, ok := OrderDirectionMap[orderDirection]
		if !ok {
			return filters, errorx.NewError("invalid order direction", errorx.BadRequest)
		}
		filters.OrderDirection = &o
	}
	return filters, nil
}
