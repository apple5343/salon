package models

import (
	"salon/internal/models"
	httpModels "salon/internal/transport/http/models"

	"github.com/brianvoe/gofakeit"
	"github.com/shopspring/decimal"
)

var PaymentTypes = []models.PaymentType{
	models.PaymentTypeCash,
	models.PaymentTypeCard,
}

func GenerateSale(car *httpModels.CarInternalResponse, clientID string) *httpModels.Sale {
	carPrice, err := decimal.NewFromString(car.Price)
	if err != nil {
		panic(err)
	}
	discount := gofakeit.Number(0, int(carPrice.IntPart())*10/100)
	discountAmount := decimal.NewFromInt(int64(discount))
	return &httpModels.Sale{
		CarID:          car.ID,
		ClientID:       clientID,
		DiscountAmount: discountAmount.String(),
		PaymentType:    string(PaymentTypes[gofakeit.Number(0, len(PaymentTypes)-1)]),
	}
}
