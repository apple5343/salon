package analytics

import (
	"salon/internal/transport/http/models"
	"time"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidTimeFormat = errorx.NewError("invalid time format", errorx.BadRequest)
)

func timeRangeFromRequest(c echo.Context) (dateFrom, dateTo *time.Time, err error) {
	var tFrom, tTo time.Time
	if dateFromStr := c.QueryParam("date_from"); dateFromStr != "" {
		tFrom, err = time.Parse(models.TimeLayout, dateFromStr)
		if err != nil {
			return
		}
		dateFrom = &tFrom
	}
	if dateToStr := c.QueryParam("date_to"); dateToStr != "" {
		tTo, err = time.Parse(models.TimeLayout, dateToStr)
		if err != nil {
			return
		}
		dateTo = &tTo
	}
	return
}
