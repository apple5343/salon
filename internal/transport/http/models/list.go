package models

import "salon/internal/models"

var DirectionMap = map[string]models.OrderDirection{
	"asc":  models.OrderDirectionASC,
	"desc": models.OrderDirectionDESC,
}
