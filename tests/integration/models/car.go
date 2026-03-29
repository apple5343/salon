package models

import (
	serviceModels "salon/internal/models"
	httpModels "salon/internal/transport/http/models"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit"
)

var CarColors = []string{
	"red",
	"blue",
	"green",
	"yellow",
	"purple",
	"orange",
	"brown",
}

func RandomVin() string {
	vin := strings.Builder{}
	vin.Grow(17)
	for i := 0; i < 17; i++ {
		n := gofakeit.Number(0, 9)
		vin.WriteString(strconv.Itoa(n))
	}
	return vin.String()
}

func RandomPrice() string {
	return strconv.Itoa(gofakeit.Number(100, 1000)) + "." + strconv.Itoa(gofakeit.Number(0, 99))
}

func GenerateCar(modelID string, supplierID string) *httpModels.Car {
	return &httpModels.Car{
		ModelID:       modelID,
		SupplierID:    supplierID,
		Vin:           RandomVin(),
		Year:          gofakeit.Number(2000, 2022),
		Color:         CarColors[gofakeit.Number(0, len(CarColors)-1)],
		InteriorColor: CarColors[gofakeit.Number(0, len(CarColors)-1)],
		Mileage:       gofakeit.Number(0, 100),
		Price:         RandomPrice(),
		Status:        string(serviceModels.CarStatusIncoming),
	}
}

func CarInternalToCar(car *httpModels.CarInternalResponse) *httpModels.Car {
	return &httpModels.Car{
		ID:            car.ID,
		ModelID:       car.Model.ID,
		SupplierID:    car.Supplier.ID,
		Vin:           car.Vin,
		Year:          car.Year,
		Color:         car.Color,
		InteriorColor: car.InteriorColor,
		Mileage:       car.Mileage,
		Price:         car.Price,
		Status:        car.Status,
		Options:       car.Options,
		CreatedAt:     car.CreatedAt,
		UpdatedAt:     car.UpdatedAt,
	}
}
