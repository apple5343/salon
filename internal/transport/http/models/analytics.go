package models

import (
	models "salon/internal/models"
)

type CarsPopularity struct {
	TransmissionTypesPopularity map[models.TransmissionType]int `json:"transmission_types_popularity"`
	FuelTypesPopularity         map[models.FuelType]int         `json:"fuel_types_popularity"`
	BodyTypesPopularity         map[models.BodyType]int         `json:"body_types_popularity"`
	DriveTypesPopularity        map[models.DriveType]int        `json:"drive_types_popularity"`
}

type SalesAnalytics struct {
	TotalDeals             int                        `json:"total_deals"`
	TotalCanceled          int                        `json:"total_canceled"`
	CanceledRevenue        string                     `json:"canceled_revenue"`
	TotalСompleted         int                        `json:"total_completed"`
	CompletedRevenue       string                     `json:"completed_revenue"`
	AverageCheck           string                     `json:"average_check"`
	ConversionRate         string                     `json:"conversion_rate"`
	DiscountImpact         string                     `json:"discount_impact"`
	AverageMargin          string                     `json:"average_margin"`
	PaymentTypesPopularity map[models.PaymentType]int `json:"payment_types_popularity"`
	CarsPopularity         CarsPopularity             `json:"cars_popularity"`
	BrandsPopularity       map[string]int             `json:"brands_popularity"`
	ModelsPopularity       map[string]int             `json:"models_popularity"`
}

func SalesAnalyticsToHttp(salesAnalytics *models.SalesAnalytics) *SalesAnalytics {
	return &SalesAnalytics{
		TotalDeals:             salesAnalytics.TotalDeals,
		TotalCanceled:          salesAnalytics.TotalCanceled,
		CanceledRevenue:        salesAnalytics.CanceledRevenue.String(),
		TotalСompleted:         salesAnalytics.TotalСompleted,
		CompletedRevenue:       salesAnalytics.CompletedRevenue.String(),
		AverageCheck:           salesAnalytics.AverageCheck.String(),
		ConversionRate:         salesAnalytics.ConversionRate.String(),
		DiscountImpact:         salesAnalytics.DiscountImpact.String(),
		AverageMargin:          salesAnalytics.AverageMargin.String(),
		PaymentTypesPopularity: salesAnalytics.PaymentTypesPopularity,
		CarsPopularity:         CarsPopularityToHttp(salesAnalytics.CarsPopularity),
		BrandsPopularity:       salesAnalytics.BrandsPopularity,
		ModelsPopularity:       salesAnalytics.ModelsPopularity,
	}
}

func CarsPopularityToHttp(carsPopularity models.CarsPopularity) CarsPopularity {
	return CarsPopularity{
		TransmissionTypesPopularity: carsPopularity.TransmissionTypesPopularity,
		FuelTypesPopularity:         carsPopularity.FuelTypesPopularity,
		BodyTypesPopularity:         carsPopularity.BodyTypesPopularity,
		DriveTypesPopularity:        carsPopularity.DriveTypesPopularity,
	}
}

type WarehouseAnalytics struct {
	TotalCars    int                      `json:"total_cars"`
	CarsByStatus map[models.CarStatus]int `json:"cars_by_status"`
	Turnover     string                   `json:"turnover"`
	CarsCount    CarsPopularity           `json:"cars_count"`
	BrandsCount  map[string]int           `json:"brands_count"`
	ModelsCount  map[string]int           `json:"models_count"`
}

func WarehouseAnalyticsToHttp(warehouseAnalytics *models.WarehouseAnalytics) *WarehouseAnalytics {
	return &WarehouseAnalytics{
		TotalCars:    warehouseAnalytics.TotalCars,
		CarsByStatus: warehouseAnalytics.CarsByStatus,
		Turnover:     warehouseAnalytics.Turnover.String(),
		CarsCount:    CarsPopularityToHttp(warehouseAnalytics.CarsCount),
		BrandsCount:  warehouseAnalytics.BrandsCount,
		ModelsCount:  warehouseAnalytics.ModelsCount,
	}
}

type SupplierAnalytics struct {
	SupplierID          string `json:"supplier_id"`
	SupplierName        string `json:"supplier_name"`
	TotalCars           int    `json:"total_cars"`
	AverageDeliveryTime int    `json:"average_delivery_time"`
	TotalPrice          string `json:"total_price"`
}

func SupplierAnalyticsToHttp(supplierAnalytics *models.SupplierAnalytics) *SupplierAnalytics {
	return &SupplierAnalytics{
		SupplierID:          supplierAnalytics.SupplierID,
		SupplierName:        supplierAnalytics.SupplierName,
		TotalCars:           supplierAnalytics.TotalCars,
		AverageDeliveryTime: supplierAnalytics.AverageDeliveryTime,
		TotalPrice:          supplierAnalytics.TotalPrice.String(),
	}
}

type SupplyAnalytics struct {
	NewCars            int            `json:"new_cars"`
	CarsOnWay          int            `json:"cars_on_way"`
	ArrivedCars        int            `json:"arrived_cars"`
	ArrivedCarsCount   CarsPopularity `json:"arrived_cars_count"`
	ArrivedBrandsCount map[string]int `json:"arrived_brands_count"`
	ArrivedModelsCount map[string]int `json:"arrived_models_count"`
}

func SupplyAnalyticsToHttp(supplyAnalytics *models.SupplyAnalytics) *SupplyAnalytics {
	return &SupplyAnalytics{
		NewCars:            supplyAnalytics.NewCars,
		CarsOnWay:          supplyAnalytics.CarsOnWay,
		ArrivedCars:        supplyAnalytics.ArrivedCars,
		ArrivedCarsCount:   CarsPopularityToHttp(supplyAnalytics.ArrivedCarsCount),
		ArrivedBrandsCount: supplyAnalytics.ArrivedBrandsCount,
		ArrivedModelsCount: supplyAnalytics.ArrivedModelsCount,
	}
}
