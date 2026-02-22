package postgres

import (
	"context"
	"database/sql"
	"fmt"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	"strings"
)

func (r *brandRepository) GetByID(ctx context.Context, id string) (*service.Brand, error) {
	var b models.Brand
	if err := r.db.GetContext(ctx, &b, "SELECT * FROM brands WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.BrandToService(&b), nil
}

func (r *brandRepository) GetBrandsByFilter(ctx context.Context, filter *service.BrandFilters) ([]*service.Brand, error) {
	var brandsList []*models.Brand
	query := `SELECT * FROM brands`
	conditions := []string{}
	args := []interface{}{}
	addCondition := func(condition string, val interface{}) {
		args = append(args, val)
		conditions = append(conditions, fmt.Sprintf("%s $%d", condition, len(args)))
	}
	if filter.Name != nil {
		addCondition("name ILIKE", "%"+*filter.Name+"%")
	}
	if filter.CountryCode != nil {
		addCondition("country_code =", *filter.CountryCode)
	}
	if len(conditions) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	}
	if filter.OrderBy != nil {
		sortMap := map[service.BrandOrderBy]string{
			service.BrandOrderByCreatedAt: "created_at",
			service.BrandOrderByUpdatedAt: "updated_at",
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
	err := r.db.SelectContext(ctx, &brandsList, query, args...)
	if err != nil {
		return nil, err
	}
	result := make([]*service.Brand, len(brandsList))
	for i, b := range brandsList {
		result[i] = models.BrandToService(b)
	}
	return result, nil
}
