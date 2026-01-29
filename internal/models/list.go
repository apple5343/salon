package models

type OrderDirection string

const (
	OrderDirectionASC  = "asc"
	OrderDirectionDESC = "desc"
)

type BaseList struct {
	Limit          *int `validate:"omitempty,min=1"`
	Offset         *int `validate:"omitempty,min=0"`
	OrderDirection *OrderDirection
}
