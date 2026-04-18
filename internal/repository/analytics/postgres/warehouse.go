package postgres

import (
	"context"
	"salon/internal/models"
	"time"

	"github.com/shopspring/decimal"
)

func (r *analyticsRepository) GetWarehouseAnalytics(ctx context.Context, turnoverDateFrom, turnoverDateTo *time.Time) (*models.WarehouseAnalytics, error) {
	analytics := &models.WarehouseAnalytics{}

	err := r.warehouseAnalyticsCarsByStatus(ctx, analytics)
	if err != nil {
		return nil, err
	}

	err = r.warehouseAnalyticsTurnover(ctx, turnoverDateFrom, turnoverDateTo, analytics)
	if err != nil {
		return nil, err
	}

	err = r.warehouseAnalyticsPopularity(ctx, analytics)
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

func (r *analyticsRepository) warehouseAnalyticsCarsByStatus(ctx context.Context, analytics *models.WarehouseAnalytics) error {
	analytics.CarsByStatus = make(map[models.CarStatus]int)

	query := `
		SELECT c.status, COUNT(*) as count
		FROM cars c
		WHERE c.status != $1
		GROUP BY c.status
	`
	rows, err := r.db.QueryContext(ctx, query, models.CarStatusSold)
	if err != nil {
		return err
	}
	defer rows.Close()

	totalCars := 0
	for rows.Next() {
		var status models.CarStatus
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return err
		}
		analytics.CarsByStatus[status] = count
		totalCars += count
	}

	allStatuses := []models.CarStatus{
		models.CarStatusAvailable,
		models.CarStatusPending,
		models.CarStatusBooked,
		models.CarStatusIncoming,
		models.CarStatusArchived,
	}
	for _, status := range allStatuses {
		if _, exists := analytics.CarsByStatus[status]; !exists {
			analytics.CarsByStatus[status] = 0
		}
	}

	analytics.TotalCars = totalCars
	return nil
}

func (r *analyticsRepository) warehouseAnalyticsTurnover(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.WarehouseAnalytics) error {
	conditions := []string{}
	args := []interface{}{}
	conditions, args = addCondition(conditions, "s.status =", args, models.SaleStatusCompleted)

	if dateFrom != nil {
		conditions, args = addCondition(conditions, "s.sale_date >=", args, dateFrom)
	}
	if dateTo != nil {
		conditions, args = addCondition(conditions, "s.sale_date <=", args, dateTo)
	}

	whereClause := getCondition(conditions)

	query := `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (s.sale_date::timestamp - c.created_at)) / 86400), 0) as avg_days
		FROM sales s
		JOIN cars c ON s.car_id = c.id
		` + whereClause

	var avgDays float64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&avgDays)
	if err != nil {
		return err
	}

	analytics.Turnover = decimal.NewFromFloat(avgDays)
	return nil
}

func (r *analyticsRepository) warehouseAnalyticsPopularity(ctx context.Context, analytics *models.WarehouseAnalytics) error {
	analytics.BrandsCount = make(map[string]int)
	analytics.ModelsCount = make(map[string]int)

	brandsQuery := `
		SELECT b.name, COUNT(*) as count
		FROM cars c
		JOIN models m ON c.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE c.status != $1
		GROUP BY b.name
		ORDER BY count DESC
	`
	rows, err := r.db.QueryContext(ctx, brandsQuery, models.CarStatusSold)
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
		analytics.BrandsCount[brandName] = count
	}

	modelsQuery := `
		SELECT m.name, COUNT(*) as count
		FROM cars c
		JOIN models m ON c.model_id = m.id
		WHERE c.status != $1
		GROUP BY m.name
		ORDER BY count DESC
	`
	rows, err = r.db.QueryContext(ctx, modelsQuery, models.CarStatusSold)
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
		analytics.ModelsCount[modelName] = count
	}

	err = r.warehouseAnalyticsCars(ctx, analytics)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) warehouseAnalyticsCars(ctx context.Context, analytics *models.WarehouseAnalytics) error {
	popularity := models.CarsPopularity{
		TransmissionTypesPopularity: make(map[models.TransmissionType]int),
		FuelTypesPopularity:         make(map[models.FuelType]int),
		BodyTypesPopularity:         make(map[models.BodyType]int),
		DriveTypesPopularity:        make(map[models.DriveType]int),
	}

	transmissionQuery := `
		SELECT m.transmission_type, COUNT(*) as count
		FROM cars c
		JOIN models m ON c.model_id = m.id
		WHERE c.status != $1
		GROUP BY m.transmission_type
	`
	rows, err := r.db.QueryContext(ctx, transmissionQuery, models.CarStatusSold)
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
		FROM cars c
		JOIN models m ON c.model_id = m.id
		WHERE c.status != $1
		GROUP BY m.fuel_type
	`
	rows, err = r.db.QueryContext(ctx, fuelQuery, models.CarStatusSold)
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
		FROM cars c
		JOIN models m ON c.model_id = m.id
		WHERE c.status != $1
		GROUP BY m.body_type
	`
	rows, err = r.db.QueryContext(ctx, bodyQuery, models.CarStatusSold)
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
		FROM cars c
		JOIN models m ON c.model_id = m.id
		WHERE c.status != $1
		GROUP BY m.drive_type
	`
	rows, err = r.db.QueryContext(ctx, driveQuery, models.CarStatusSold)
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

	analytics.CarsCount = popularity
	return nil
}
