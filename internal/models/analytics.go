package models

import "github.com/shopspring/decimal"

// Статистика расчитывается за определенный период

type CarsPopularity struct {
	TransmissionTypesPopularity map[TransmissionType]int
	FuelTypesPopularity         map[FuelType]int
	BodyTypesPopularity         map[BodyType]int
	DriveTypesPopularity        map[DriveType]int
}

type SalesAnalytics struct {
	TotalRevenue           decimal.Decimal // Сумма всех final_price.
	TotalSales             int
	AverageCheck           decimal.Decimal
	ConversionRate         decimal.Decimal     // (SaleCompleted) ко всем начатым (SaleCreated)
	DiscountImpact         decimal.Decimal     // Сумма всех discount_amount
	AverageMargin          decimal.Decimal     // Средняя разница между final_price и origin_price
	PaymentTypesPopularity map[PaymentType]int // Количество продаж по типу оплаты
	CarsPopularity         CarsPopularity      // Количество проданных машин по параметрам
	BrandsPopularity       map[string]int      // Количество проданных машин по брендам
	ModelsPopularity       map[string]int      // Количество проданных машин по моделям
}

type WarehouseAnalytics struct { // В учет не берутся проданные машины, статиститка не зависит от периода
	CarsByStatus map[CarStatus]int // Количество машин по статусу (кроме проданных)
	TotalCars    int
	Turnover     decimal.Decimal // Разница во времени между cars.created_at (поступление) и sales.sale_date (продажа). Брать среднюю?
	CarsCount    CarsPopularity  // Количество машин по параметрам
	BrandsCount  map[string]int  // Количество машин по брендам
	ModelsCount  map[string]int  // Количество машин по моделям
}

type EmployeeAnalytics struct {
	EmployeeID      string
	EmployeeName    string
	TotalSales      int
	TotalPrice      decimal.Decimal
	AverageDiscount decimal.Decimal
	AverageCheck    decimal.Decimal
}

type ModelAnalytics struct {
	ModelID   string
	ModelName string
	// Не знаю, какую статистику сюда добавить
}

type SupplierAnalytics struct {
	SupplierID          string
	SupplierName        string
	TotalCars           int             // Количество поставленных машин
	AverageDeliveryTime int             // Среднее время доставки машины
	TotalPrice          decimal.Decimal // Сумма поставленных машин
}

type SupplyAnalytics struct {
	NewCars     int // Количество машин, которые были созданы
	ArrivedCars int // Количество машин, которые были поставлены
	CarsCount   CarsPopularity
	BrandsCount map[string]int
	ModelsCount map[string]int
}
