package models

import "salon/internal/models"

var OrderDirectionMap = map[string]models.OrderDirection{
	"asc":  models.OrderDirectionASC,
	"desc": models.OrderDirectionDESC,
}
