package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	"strings"

	service "salon/internal/models"
)

func (r *modelRepository) GetByID(ctx context.Context, id string) (*service.Model, error) {
	var m models.Model
	if err := r.db.GetContext(ctx, &m, "SELECT * FROM models WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.ModelToService(&m), nil
}

func (r *modelRepository) GetModelsByFilter(ctx context.Context, filter *service.ModelFilters) ([]*service.ModelShort, error) {
	var modelsList []*models.ModelShort
	query := `SELECT m.id, b.name as brand_name, m.name, m.generation, m.body_type, m.drive_type, m.power_hp, m.base_price
		FROM models m
		JOIN brands b ON m.brand_id = b.id`
	conditions := []string{}
	args := []interface{}{}
	addCondition := func(condition string, val interface{}) {
		args = append(args, val)
		conditions = append(conditions, fmt.Sprintf("%s $%d", condition, len(args)))
	}
	if filter.BrandID != nil {
		addCondition("b.id =", *filter.BrandID)
	}
	if filter.BodyType != nil {
		addCondition("m.body_type =", *filter.BodyType)
	}
	if filter.DriveType != nil {
		addCondition("m.drive_type =", *filter.DriveType)
	}
	if filter.MinEngineDisplacement != nil {
		addCondition("m.engine_displacement >=", *filter.MinEngineDisplacement)
	}
	if filter.MaxEngineDisplacement != nil {
		addCondition("m.engine_displacement <=", *filter.MaxEngineDisplacement)
	}
	if filter.MinPowerHP != nil {
		addCondition("m.power_hp >=", *filter.MinPowerHP)
	}
	if filter.MaxPowerHP != nil {
		addCondition("m.power_hp <=", *filter.MaxPowerHP)
	}
	if filter.MinBasePrice != nil {
		addCondition("m.base_price >=", *filter.MinBasePrice)
	}
	if filter.MaxBasePrice != nil {
		addCondition("m.base_price <=", *filter.MaxBasePrice)
	}
	if filter.Name != nil {
		addCondition("m.name ILIKE", "%"+*filter.Name+"%")
	}
	if filter.Generation != nil {
		addCondition("m.generation ILIKE", "%"+*filter.Generation+"%")
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if filter.OrderBy != nil {
		sortMap := map[service.ModelOrderBy]string{
			service.ModelOrderByBasePrice:          "m.base_price",
			service.ModelOrderByEngineDisplacement: "m.engine_displacement",
			service.ModelOrderByPowerHP:            "m.power_hp",
			service.ModelOrderByName:               "m.name",
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
	err := r.db.SelectContext(ctx, &modelsList, query, args...)
	if err != nil {
		return nil, err
	}
	result := make([]*service.ModelShort, len(modelsList))
	for i, m := range modelsList {
		result[i] = models.ModelShortToService(m)
	}
	return result, nil
}
