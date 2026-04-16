package generator

import (
	"math/rand/v2"
	"salon/internal/models"

	"github.com/shopspring/decimal"
)

func (g *Generator) GenerateSale(car *models.Car, clientID string, employeeID string) *models.Sale {
	discount := rand.IntN(int(car.Price.IntPart()) * 10 / 100)
	paymentTypes := []models.PaymentType{models.PaymentTypeCash, models.PaymentTypeCard}
	paymentType := paymentTypes[rand.IntN(len(paymentTypes))]
	return &models.Sale{
		CarID:          car.ID,
		ClientID:       clientID,
		EmployeeID:     employeeID,
		DiscountAmount: decimal.NewFromInt(int64(discount)),
		PaymentType:    paymentType,
	}
}
