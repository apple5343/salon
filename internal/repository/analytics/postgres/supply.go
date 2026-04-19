package postgres

import (
	"context"
	"fmt"
	"salon/internal/models"
	"time"

	"github.com/lib/pq"
)

func (r *analyticsRepository) GetSupplyAnalytics(ctx context.Context, dateFrom, dateTo *time.Time) (*models.SupplyAnalytics, error) {
	analytics := &models.SupplyAnalytics{}

	err := r.supplyAnalyticsNewCars(ctx, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}

	err = r.supplyAnalyticsCarsOnWay(ctx, analytics)
	if err != nil {
		return nil, err
	}

	err = r.supplyAnalyticsArrivedCars(ctx, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}

	err = r.supplyAnalyticsArrivedPopularity(ctx, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

func (r *analyticsRepository) supplyAnalyticsNewCars(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.SupplyAnalytics) error {
	conditions := []string{}
	args := []interface{}{}

	if dateFrom != nil {
		conditions, args = addCondition(conditions, "created_at >=", args, dateFrom)
	}
	if dateTo != nil {
		conditions, args = addCondition(conditions, "created_at <=", args, dateTo)
	}

	whereClause := getCondition(conditions)

	query := `SELECT COUNT(*) FROM cars ` + whereClause

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&analytics.NewCars)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) supplyAnalyticsCarsOnWay(ctx context.Context, analytics *models.SupplyAnalytics) error {
	query := `SELECT COUNT(*) FROM cars WHERE status = $1`

	err := r.db.QueryRowContext(ctx, query, models.CarStatusIncoming).Scan(&analytics.CarsOnWay)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) supplyAnalyticsArrivedCars(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.SupplyAnalytics) error {
	conditions := []string{}
	args := []interface{}{}

	if dateFrom != nil {
		conditions, args = addCondition(conditions, "e.created_at >=", args, dateFrom)
	}
	if dateTo != nil {
		conditions, args = addCondition(conditions, "e.created_at <=", args, dateTo)
	}
	conditions, args = addCondition(conditions, "1 =", args, 1)

	whereClause := getCondition(conditions)

	query := `
		WITH status_changes AS (
			SELECT DISTINCT e.entity_id
			FROM events e
			` + whereClause + `
				AND e.event_type = $` + fmt.Sprintf("%d", len(args)+1) + `
				AND e.entity_type = $` + fmt.Sprintf("%d", len(args)+2) + `
				AND (e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+3) + ` OR e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+4) + `)
		),
		incoming_cars AS (
			SELECT DISTINCT e.entity_id
			FROM events e
			WHERE e.event_type = $` + fmt.Sprintf("%d", len(args)+5) + `
				AND e.entity_type = $` + fmt.Sprintf("%d", len(args)+6) + `
				AND e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+7) + `
		)
		SELECT COUNT(*)
		FROM status_changes sc
		JOIN incoming_cars ic ON sc.entity_id = ic.entity_id
	`

	args = append(args,
		models.EventTypeUpdated,
		models.EntityTypeCar,
		models.CarStatusPending,
		models.CarStatusAvailable,
		models.EventTypeCreated,
		models.EntityTypeCar,
		models.CarStatusIncoming,
	)

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&analytics.ArrivedCars)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) supplyAnalyticsArrivedPopularity(ctx context.Context, dateFrom, dateTo *time.Time, analytics *models.SupplyAnalytics) error {
	analytics.ArrivedBrandsCount = make(map[string]int)
	analytics.ArrivedModelsCount = make(map[string]int)

	conditions := []string{}
	args := []interface{}{}

	if dateFrom != nil {
		conditions, args = addCondition(conditions, "e.created_at >=", args, dateFrom)
	}
	if dateTo != nil {
		conditions, args = addCondition(conditions, "e.created_at <=", args, dateTo)
	}
	conditions, args = addCondition(conditions, "1 =", args, 1)

	whereClause := getCondition(conditions)

	arrivedCarsQuery := `
		WITH status_changes AS (
			SELECT DISTINCT e.entity_id
			FROM events e
			` + whereClause + `
				AND e.event_type = $` + fmt.Sprintf("%d", len(args)+1) + `
				AND e.entity_type = $` + fmt.Sprintf("%d", len(args)+2) + `
				AND (e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+3) + ` OR e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+4) + `)
		),
		incoming_cars AS (
			SELECT DISTINCT e.entity_id
			FROM events e
			WHERE e.event_type = $` + fmt.Sprintf("%d", len(args)+5) + `
				AND e.entity_type = $` + fmt.Sprintf("%d", len(args)+6) + `
				AND e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+7) + `
		)
		SELECT sc.entity_id
		FROM status_changes sc
		JOIN incoming_cars ic ON sc.entity_id = ic.entity_id
	`

	arrivedArgs := append([]interface{}{}, args...)
	arrivedArgs = append(arrivedArgs,
		models.EventTypeUpdated,
		models.EntityTypeCar,
		models.CarStatusPending,
		models.CarStatusAvailable,
		models.EventTypeCreated,
		models.EntityTypeCar,
		models.CarStatusIncoming,
	)

	rows, err := r.db.QueryContext(ctx, arrivedCarsQuery, arrivedArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var arrivedCarIDs []string
	for rows.Next() {
		var carID string
		if err := rows.Scan(&carID); err != nil {
			return err
		}
		arrivedCarIDs = append(arrivedCarIDs, carID)
	}

	if len(arrivedCarIDs) == 0 {
		analytics.ArrivedCarsCount = models.CarsPopularity{
			TransmissionTypesPopularity: make(map[models.TransmissionType]int),
			FuelTypesPopularity:         make(map[models.FuelType]int),
			BodyTypesPopularity:         make(map[models.BodyType]int),
			DriveTypesPopularity:        make(map[models.DriveType]int),
		}
		return nil
	}

	brandsQuery := `
		SELECT b.name, COUNT(*) as count
		FROM cars c
		JOIN models m ON c.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE c.id = ANY($1)
		GROUP BY b.name
		ORDER BY count DESC
	`
	rows, err = r.db.QueryContext(ctx, brandsQuery, pq.Array(arrivedCarIDs))
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
		analytics.ArrivedBrandsCount[brandName] = count
	}

	modelsQuery := `
		SELECT m.name, COUNT(*) as count
		FROM cars c
		JOIN models m ON c.model_id = m.id
		WHERE c.id = ANY($1)
		GROUP BY m.name
		ORDER BY count DESC
	`
	rows, err = r.db.QueryContext(ctx, modelsQuery, pq.Array(arrivedCarIDs))
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
		analytics.ArrivedModelsCount[modelName] = count
	}

	err = r.supplyAnalyticsArrivedCarsCount(ctx, arrivedCarIDs, analytics)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) supplyAnalyticsArrivedCarsCount(ctx context.Context, carIDs []string, analytics *models.SupplyAnalytics) error {
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
		WHERE c.id = ANY($1)
		GROUP BY m.transmission_type
	`
	rows, err := r.db.QueryContext(ctx, transmissionQuery, pq.Array(carIDs))
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
		WHERE c.id = ANY($1)
		GROUP BY m.fuel_type
	`
	rows, err = r.db.QueryContext(ctx, fuelQuery, pq.Array(carIDs))
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
		WHERE c.id = ANY($1)
		GROUP BY m.body_type
	`
	rows, err = r.db.QueryContext(ctx, bodyQuery, pq.Array(carIDs))
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
		WHERE c.id = ANY($1)
		GROUP BY m.drive_type
	`
	rows, err = r.db.QueryContext(ctx, driveQuery, pq.Array(carIDs))
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

	analytics.ArrivedCarsCount = popularity
	return nil
}
