package postgres

import (
	"context"
	"fmt"
	service "salon/internal/models"
	"salon/internal/repository/models"
	"strings"
)

func (r *employeeRepository) GetEmployeesByFilter(ctx context.Context, filter *service.EmployeeFilters) ([]*service.Employee, error) {
	var employeesList []*models.Employee
	query := `SELECT * FROM employees`
	conditions := []string{}
	args := []interface{}{}
	addCondition := func(condition string, val interface{}) {
		args = append(args, val)
		conditions = append(conditions, fmt.Sprintf("%s $%d", condition, len(args)))
	}
	if filter.FullName != nil {
		addCondition("full_name ILIKE", "%"+*filter.FullName+"%")
	}
	if filter.Role != nil {
		addCondition("role =", *filter.Role)
	}
	if filter.Status != nil {
		addCondition("status =", *filter.Status)
	}
	if len(conditions) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	}
	if filter.OrderBy != nil {
		sortMap := map[service.EmployeeOrderBy]string{
			service.EmployeeOrderByCreatedAt: "created_at",
			service.EmployeeOrderByHireDate:  "hire_date",
		}
		direction := "ASC"
		if filter.OrderDirection != nil && *filter.OrderDirection == service.OrderDirectionDESC {
			direction = "DESC"
		}
		if dbField, ok := sortMap[*filter.OrderBy]; ok {
			query += fmt.Sprintf(" ORDER BY %s %s NULLS LAST, %s %s", dbField, direction, "created_at", direction)
		}
	} else {
		query += " ORDER BY created_at ASC"
	}
	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d ", len(args)+1)
		args = append(args, *filter.Limit)
	}
	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, *filter.Offset)
	}
	err := r.db.SelectContext(ctx, &employeesList, query, args...)
	if err != nil {
		return nil, err
	}
	result := make([]*service.Employee, len(employeesList))
	for i, e := range employeesList {
		result[i] = models.EmployeeToService(e)
	}
	return result, nil
}
