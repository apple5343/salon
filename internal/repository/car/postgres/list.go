package postgres

import (
	"context"
	"fmt"
	"strings"

	service "salon/internal/models"
	"salon/internal/repository/models"
)

func (r *carRepository) GetCarsByFilter(ctx context.Context, filter *service.CarFilters) ([]*service.CarShort, error) {
	var cars []*models.CarShort
	query := `SELECT c.id, m.name as model_name, s.name as supplier_name, b.name as brand_name, c.vin, c.status, c.price, c.year
		FROM cars c
		JOIN models m ON c.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		JOIN suppliers s ON c.supplier_id = s.id`
	conditions := []string{}
	args := []interface{}{}
	addCondition := func(condition string, val interface{}) {
		args = append(args, val)
		conditions = append(conditions, fmt.Sprintf("%s $%d", condition, len(args)))
	}
	if filter.SupplierID != nil {
		addCondition("s.id =", *filter.SupplierID)
	}
	if filter.ModelID != nil {
		addCondition("m.id =", *filter.ModelID)
	}
	if filter.BrandID != nil {
		addCondition("b.id =", *filter.BrandID)
	}
	if filter.Color != nil {
		addCondition("c.color =", *filter.Color)
	}
	if filter.InteriorColor != nil {
		addCondition("c.interior_color =", *filter.InteriorColor)
	}
	if filter.Status != nil {
		addCondition("c.status =", *filter.Status)
	}
	if filter.MinPrice != nil {
		addCondition("c.price >=", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		addCondition("c.price <=", *filter.MaxPrice)
	}
	if filter.MinYear != nil {
		addCondition("c.year >=", *filter.MinYear)
	}
	if filter.MaxYear != nil {
		addCondition("c.year <=", *filter.MaxYear)
	}
	if filter.MinMileage != nil {
		addCondition("c.mileage >=", *filter.MinMileage)
	}
	if filter.MaxMileage != nil {
		addCondition("c.mileage <=", *filter.MaxMileage)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if filter.OrderBy != nil {
		sortMap := map[service.CarOrderBy]string{
			service.CarOrderByPrice:     "price",
			service.CarOrderByYear:      "year",
			service.CarOrderByMileage:   "mileage",
			service.CarOrderByCreatedAt: "c.created_at",
			service.CarOrderByUpdatedAt: "c.updated_at",
		}
		direction := "ASC"
		if filter.OrderDirection != nil && *filter.OrderDirection == service.OrderDirectionDESC {
			direction = "DESC"
		}
		if dbField, ok := sortMap[*filter.OrderBy]; ok {
			query += fmt.Sprintf(" ORDER BY %s %s", dbField, direction)
		}
	}
	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d ", len(args)+1)
		args = append(args, *filter.Limit)
	}
	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, *filter.Offset)
	}
	err := r.db.SelectContext(ctx, &cars, query, args...)
	if err != nil {
		return nil, err
	}
	result := make([]*service.CarShort, len(cars))
	for i, c := range cars {
		result[i] = models.CarShortToService(c)
	}
	return result, nil
}