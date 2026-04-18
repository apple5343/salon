package postgres

import (
	"context"
	"fmt"
	"salon/internal/models"
	"strings"
	"time"
)

func (r *analyticsRepository) GetSupplierAnalytics(ctx context.Context, supplierID string, dateFrom, dateTo *time.Time) (*models.SupplierAnalytics, error) {
	analytics := &models.SupplierAnalytics{
		SupplierID: supplierID,
	}

	err := r.db.QueryRowContext(ctx, "SELECT name FROM suppliers WHERE id = $1", supplierID).Scan(&analytics.SupplierName)
	if err != nil {
		return nil, err
	}

	err = r.supplierAnalyticsBasic(ctx, supplierID, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}

	err = r.supplierAnalyticsDeliveryTime(ctx, supplierID, dateFrom, dateTo, analytics)
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

func (r *analyticsRepository) supplierAnalyticsBasic(ctx context.Context, supplierID string, dateFrom, dateTo *time.Time, analytics *models.SupplierAnalytics) error {
	conditions := []string{"c.supplier_id = $1"}
	args := []interface{}{supplierID}

	if dateFrom != nil {
		conditions, args = addCondition(conditions, "c.created_at >=", args, dateFrom)
	}
	if dateTo != nil {
		conditions, args = addCondition(conditions, "c.created_at <=", args, dateTo)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	query := `
		SELECT
			COUNT(*) as total_cars,
			COALESCE(SUM(c.price), 0) as total_price
		FROM cars c
		` + whereClause

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&analytics.TotalCars,
		&analytics.TotalPrice,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) supplierAnalyticsDeliveryTime(ctx context.Context, supplierID string, dateFrom, dateTo *time.Time, analytics *models.SupplierAnalytics) error {
	conditions := []string{"c.supplier_id = $1"}
	args := []interface{}{supplierID}

	if dateFrom != nil {
		conditions, args = addCondition(conditions, "c.created_at >=", args, dateFrom)
	}
	if dateTo != nil {
		conditions, args = addCondition(conditions, "c.created_at <=", args, dateTo)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	query := `
		WITH car_created AS (
			SELECT
				e.entity_id as car_id,
				e.created_at as created_time
			FROM events e
			JOIN cars c ON e.entity_id = c.id
			` + whereClause + `
				AND e.event_type = $` + fmt.Sprintf("%d", len(args)+1) + `
				AND e.entity_type = $` + fmt.Sprintf("%d", len(args)+2) + `
				AND e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+3) + `
		),
		first_update AS (
			SELECT DISTINCT ON (e.entity_id)
				e.entity_id as car_id,
				e.created_at as updated_time
			FROM events e
			WHERE e.event_type = $` + fmt.Sprintf("%d", len(args)+4) + `
				AND e.entity_type = $` + fmt.Sprintf("%d", len(args)+5) + `
				AND (e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+6) + ` OR e.payload->>'status' = $` + fmt.Sprintf("%d", len(args)+7) + `)
				AND e.entity_id IN (SELECT car_id FROM car_created)
			ORDER BY e.entity_id, e.created_at ASC
		)
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (fu.updated_time - cc.created_time)) / 86400), 0) as avg_days
		FROM car_created cc
		JOIN first_update fu ON cc.car_id = fu.car_id
	`

	args = append(args,
		models.EventTypeCreated,
		models.EntityTypeCar,
		models.CarStatusIncoming,
		models.EventTypeUpdated,
		models.EntityTypeCar,
		models.CarStatusPending,
		models.CarStatusAvailable,
	)

	var avgDays float64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&avgDays)
	if err != nil {
		return err
	}

	analytics.AverageDeliveryTime = int(avgDays)
	return nil
}
