package postgres

import (
	"context"
	"fmt"
	service "salon/internal/models"
	"salon/internal/repository/models"
	"strings"
)

func (r *supplierRepository) GetSuppliersByFilter(ctx context.Context, filter *service.SupplierFilters) ([]*service.Supplier, error) {
	var suppliers []*models.Supplier

	query := "SELECT * FROM suppliers"
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
		sortMap := map[service.SupplierOrderBy]string{
			service.SupplierOrderByCreatedAt: "created_at",
			service.SupplierOrderByUpdatedAt: "updated_at",
		}
		direcion := "ASC"
		if filter.OrderDirection != nil {
			direcion = string(*filter.OrderDirection)
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortMap[*filter.OrderBy], direcion)
	}
	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}
	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}
	err := r.db.SelectContext(ctx, &suppliers, query, args...)
	if err != nil {
		return nil, err
	}
	result := make([]*service.Supplier, len(suppliers))
	for i, s := range suppliers {
		result[i] = models.SupplierToService(s)
	}
	return result, nil
}
