package postgres

import (
	"context"
	"salon/internal/models"
	"strings"
	"time"
)

func (r *analyticsRepository) GetEmployeeAnalytics(ctx context.Context, employeeID string, dateFrom, dateTo *time.Time) (*models.EmployeeAnalytics, error) {
	conditions := []string{"s.employee_id = $1", "s.status = $2"}
	args := []interface{}{employeeID, models.SaleStatusCompleted}

	if dateFrom != nil {
		conditions, args = addCondition(conditions, "s.sale_date >=", args, dateFrom)
	}
	if dateTo != nil {
		conditions, args = addCondition(conditions, "s.sale_date <=", args, dateTo)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	query := `
		SELECT
			COUNT(*) as total_sales,
			COALESCE(SUM(s.final_price), 0) as total_revenue,
			COALESCE(AVG(s.discount_amount), 0) as average_discount,
			COALESCE(AVG(s.final_price), 0) as average_check
		FROM sales s
		` + whereClause

	var employeeAnalytics models.EmployeeAnalytics
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&employeeAnalytics.TotalSales,
		&employeeAnalytics.TotalRevenue,
		&employeeAnalytics.AverageDiscount,
		&employeeAnalytics.AverageCheck,
	)
	if err != nil {
		return nil, err
	}

	var employeeName string
	err = r.db.QueryRowContext(ctx, "SELECT full_name FROM employees WHERE id = $1", employeeID).Scan(&employeeName)
	if err != nil {
		return nil, err
	}

	employeeAnalytics.EmployeeName = employeeName
	employeeAnalytics.EmployeeID = employeeID

	return &employeeAnalytics, nil
}
