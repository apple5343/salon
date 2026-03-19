package models

import (
	serviceModels "salon/internal/models"
	httpModels "salon/internal/transport/http/models"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit"
)

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
	return strconv.Itoa(gofakeit.Number(1000, 100000)) + "." + strconv.Itoa(gofakeit.Number(0, 99))
}

func GenerateCar(modelID string, supplierID string) *httpModels.Car {
	return &httpModels.Car{
		ModelID:       modelID,
		SupplierID:    supplierID,
		Vin:           RandomVin(),
		Year:          gofakeit.Number(2000, 2022),
		Color:         gofakeit.Color(),
		InteriorColor: gofakeit.Color(),
		Mileage:       gofakeit.Number(0, 100),
		Price:         RandomPrice(),
		Status:        string(serviceModels.CarStatusIncoming),
	}
}
