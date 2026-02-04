package postgres

import (
	"context"
	"fmt"
	service "salon/internal/models"
	"salon/internal/repository/models"
	"strings"
)

func (r *saleRepository) GetSalesByFilter(ctx context.Context, filter *service.SaleFilters) ([]*service.Sale, error) {
	var sales []*models.Sale

	query := "SELECT * FROM sales "
	conditions := []string{}
	args := []interface{}{}
	addCondition := func(condition string, val interface{}) {
		args = append(args, val)
		conditions = append(conditions, fmt.Sprintf("%s $%d", condition, len(args)))
	}
	if filter.CarID != nil {
		addCondition("car_id =", *filter.CarID)
	}
	if filter.ClientID != nil {
		addCondition("client_id =", *filter.ClientID)
	}
	if filter.EmployeeID != nil {
		addCondition("employee_id =", *filter.EmployeeID)
	}
	if filter.Status != nil {
		addCondition("status =", *filter.Status)
	}
	if filter.PaymentType != nil {
		addCondition("payment_type =", string(*filter.PaymentType))
	}
	if filter.DateFrom != nil {
		addCondition("sale_date >=", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		addCondition("sale_date <=", *filter.DateTo)
	}
	if filter.FinalPriceMin != nil {
		addCondition("final_price >=", *filter.FinalPriceMin)
	}
	if filter.FinalPriceMax != nil {
		addCondition("final_price <=", *filter.FinalPriceMax)
	}
	if len(conditions) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	}
	if filter.OrderBy != nil {
		sortMap := map[service.SaleOrderBy]string{
			service.SaleOrderByDate:       "sale_date",
			service.SaleOrderByFinalPrice: "final_price",
		}
		direction := "ASC"
		if filter.OrderDirection != nil {
			direction = string(*filter.OrderDirection)
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortMap[*filter.OrderBy], direction)
	}
	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}
	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}
	err := r.db.SelectContext(ctx, &sales, query, args...)
	if err != nil {
		return nil, err
	}
	result := make([]*service.Sale, len(sales))
	for i, s := range sales {
		result[i] = models.SaleToService(s)
	}
	return result, nil
}
