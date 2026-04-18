package postgres

import (
	"context"
	"fmt"
	"salon/internal/models"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func (r *analyticsRepository) GetSalesAnalytics(ctx context.Context, dateFrom, dateTo *time.Time) (*models.SalesAnalytics, error) {
	analytics := &models.SalesAnalytics{}
	err := r.salesAnalyticsCanceled(ctx, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}
	err = r.salesAnalyticsCompleted(ctx, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}
	err = r.salesAnalyticsPaymentTypes(ctx, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}
	err = r.salesAnalyticsPopularity(ctx, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}
	analytics.TotalDeals = analytics.TotalCanceled + analytics.TotalСompleted
	return analytics, nil
}

func (r *analyticsRepository) salesAnalyticsCanceled(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.SalesAnalytics) error {
	conditions := []string{}
	args := []interface{}{}
	conditions, args = addCondition(conditions, "s.status =", args, models.SaleStatusCanceled)
	conditions, args = addTimeRangeCondition(conditions, args, dateFrom, dateTo)

	whereClause := getCondition(conditions)
	query := `
		SELECT
			COALESCE(SUM(s.final_price), 0) as total_revenue,
			COUNT(s.id) as total_sales
		FROM sales s
		` + whereClause

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&analytics.CanceledRevenue,
		&analytics.TotalCanceled,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *analyticsRepository) salesAnalyticsCompleted(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.SalesAnalytics) error {
	conditions := []string{}
	args := []interface{}{}
	conditions, args = addCondition(conditions, "s.status =", args, models.SaleStatusCompleted)
	conditions, args = addTimeRangeCondition(conditions, args, dateFrom, dateTo)

	whereClause := getCondition(conditions)
	query := `
		SELECT
			COALESCE(SUM(s.final_price), 0) as total_revenue,
			COUNT(s.id) as total_sales,
			COALESCE(AVG(s.final_price), 0) as average_check,
			COALESCE(SUM(s.discount_amount), 0) as discount_impact,
			COALESCE(AVG(s.original_price - s.final_price), 0) as average_margin
		FROM sales s
		` + whereClause

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&analytics.CompletedRevenue,
		&analytics.TotalСompleted,
		&analytics.AverageCheck,
		&analytics.DiscountImpact,
		&analytics.AverageMargin,
	)
	if err != nil {
		return err
	}

	conversionConditions := []string{}
	conversionArgs := []interface{}{models.SaleStatusCompleted}
	conversionConditions, conversionArgs = addTimeRangeCondition(conversionConditions, conversionArgs, dateFrom, dateTo)
	conversionWhereClause := getCondition(conversionConditions)
	conversionQuery := `
		SELECT
			COUNT(CASE WHEN status = $1 THEN 1 END)::decimal / NULLIF(COUNT(*)::decimal, 0) as conversion_rate
		FROM sales s 
		` + conversionWhereClause

	var conversionRate decimal.NullDecimal
	err = r.db.QueryRowContext(ctx, conversionQuery, conversionArgs...).Scan(&conversionRate)
	if err != nil {
		return err
	}
	if conversionRate.Valid {
		analytics.ConversionRate = conversionRate.Decimal
	} else {
		analytics.ConversionRate = decimal.Zero
	}

	return nil
}

func (r *analyticsRepository) salesAnalyticsPopularity(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.SalesAnalytics) error {
	analytics.BrandsPopularity = make(map[string]int)
	analytics.ModelsPopularity = make(map[string]int)

	conditions := []string{}
	args := []interface{}{}
	conditions, args = addTimeRangeCondition(conditions, args, dateFrom, dateTo)
	conditions, args = addCondition(conditions, "s.status =", args, models.SaleStatusCompleted)

	whereClause := getCondition(conditions)
	brandsQuery := `
		SELECT b.name, COUNT(*) as count
		FROM sales s
		JOIN cars c ON s.car_id = c.id
		JOIN models m ON c.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		` + whereClause + `
		GROUP BY b.name
		ORDER BY count DESC
	`
	rows, err := r.db.QueryContext(ctx, brandsQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var brandName string
		var count int
		if err := rows.Scan(&brandName, &count); err != nil {
			return err
		}
		analytics.BrandsPopularity[brandName] = count
	}

	modelsQuery := `
		SELECT m.name, COUNT(*) as count
		FROM sales s
		JOIN cars c ON s.car_id = c.id
		JOIN models m ON c.model_id = m.id
		` + whereClause + `
		GROUP BY m.name
		ORDER BY count DESC
	`
	rows, err = r.db.QueryContext(ctx, modelsQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var modelName string
		var count int
		if err := rows.Scan(&modelName, &count); err != nil {
			return err
		}
		analytics.ModelsPopularity[modelName] = count
	}

	err = r.salesAnalyticsCars(ctx, dateFrom, dateTo, analytics)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) salesAnalyticsPaymentTypes(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.SalesAnalytics) error {
	analytics.PaymentTypesPopularity = make(map[models.PaymentType]int)

	conditions := []string{}
	args := []interface{}{}
	conditions, args = addCondition(conditions, "s.status =", args, models.SaleStatusCompleted)
	conditions, args = addTimeRangeCondition(conditions, args, dateFrom, dateTo)

	whereClause := getCondition(conditions)
	query := `
		SELECT s.payment_type, COUNT(*) as count
		FROM sales s
		` + whereClause + `
		GROUP BY s.payment_type
	`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var paymentType models.PaymentType
		var count int
		if err := rows.Scan(&paymentType, &count); err != nil {
			return err
		}
		analytics.PaymentTypesPopularity[paymentType] = count
	}

	return nil
}

func (r *analyticsRepository) salesAnalyticsCars(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.SalesAnalytics) error {
	popularity := models.CarsPopularity{
		TransmissionTypesPopularity: make(map[models.TransmissionType]int),
		FuelTypesPopularity:         make(map[models.FuelType]int),
		BodyTypesPopularity:         make(map[models.BodyType]int),
		DriveTypesPopularity:        make(map[models.DriveType]int),
	}

	conditions := []string{}
	args := []interface{}{}
	conditions, args = addCondition(conditions, "s.status =", args, models.SaleStatusCompleted)
	conditions, args = addTimeRangeCondition(conditions, args, dateFrom, dateTo)

	whereClause := getCondition(conditions)

	transmissionQuery := `
		SELECT m.transmission_type, COUNT(*) as count
		FROM sales s
		JOIN cars c ON s.car_id = c.id
		JOIN models m ON c.model_id = m.id
		` + whereClause + `
		GROUP BY m.transmission_type
	`
	rows, err := r.db.QueryContext(ctx, transmissionQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var transmissionType models.TransmissionType
		var count int
		if err := rows.Scan(&transmissionType, &count); err != nil {
			return err
		}
		popularity.TransmissionTypesPopularity[transmissionType] = count
	}

	fuelQuery := `
		SELECT m.fuel_type, COUNT(*) as count
		FROM sales s
		JOIN cars c ON s.car_id = c.id
		JOIN models m ON c.model_id = m.id
		` + whereClause + `
		GROUP BY m.fuel_type
	`
	rows, err = r.db.QueryContext(ctx, fuelQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var fuelType models.FuelType
		var count int
		if err := rows.Scan(&fuelType, &count); err != nil {
			return err
		}
		popularity.FuelTypesPopularity[fuelType] = count
	}

	bodyQuery := `
		SELECT m.body_type, COUNT(*) as count
		FROM sales s
		JOIN cars c ON s.car_id = c.id
		JOIN models m ON c.model_id = m.id
		` + whereClause + `
		GROUP BY m.body_type
	`
	rows, err = r.db.QueryContext(ctx, bodyQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var bodyType models.BodyType
		var count int
		if err := rows.Scan(&bodyType, &count); err != nil {
			return err
		}
		popularity.BodyTypesPopularity[bodyType] = count
	}

	driveQuery := `
		SELECT m.drive_type, COUNT(*) as count
		FROM sales s
		JOIN cars c ON s.car_id = c.id
		JOIN models m ON c.model_id = m.id
		` + whereClause + `
		GROUP BY m.drive_type
	`
	rows, err = r.db.QueryContext(ctx, driveQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var driveType models.DriveType
		var count int
		if err := rows.Scan(&driveType, &count); err != nil {
			return err
		}
		popularity.DriveTypesPopularity[driveType] = count
	}

	analytics.CarsPopularity = popularity
	return nil
}

func addCondition(conditions []string, condition string, args []interface{}, val interface{}) ([]string, []interface{}) {
	args = append(args, val)
	conditions = append(conditions, fmt.Sprintf("%s $%d", condition, len(args)))
	return conditions, args
}

func getCondition(conditions []string) string {
	if len(conditions) > 0 {
		return fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	}
	return ""
}

func addTimeRangeCondition(conditions []string, args []interface{}, dateFrom, dateTo *time.Time) ([]string, []interface{}) {
	if dateFrom != nil {
		conditions, args = addCondition(conditions, "s.created_at >=", args, dateFrom)
	}
	if dateTo != nil {
		conditions, args = addCondition(conditions, "s.created_at <=", args, dateTo)
	}
	return conditions, args
}
