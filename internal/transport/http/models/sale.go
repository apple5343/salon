package models

import (
	"salon/internal/models"
	"time"

	"github.com/apple5343/errorx"
	"github.com/shopspring/decimal"
)

var PaymentTypeMap = map[string]models.PaymentType{
	"cash": models.PaymentTypeCash,
	"card": models.PaymentTypeCard,
}

var SaleStatusMap = map[string]models.SaleStatus{
	"pending":   models.SaleStatusPending,
	"completed": models.SaleStatusCompleted,
	"canceled":  models.SaleStatusCanceled,
}

var SaleOrderByMap = map[string]models.SaleOrderBy{
	"final_price": models.SaleOrderByFinalPrice,
	"date":        models.SaleOrderByDate,
}

type Sale struct {
	ID              string    `json:"id"`
	CarID           string    `json:"car_id"`
	ClientID        string    `json:"client_id"`
	EmployeeID      string    `json:"employee_id"`
	SaleDate        string    `json:"sale_date"`
	OriginPrice     string    `json:"origin_price"`
	DiscountAmount  string    `json:"discount_amount"`
	DiscountPercent string    `json:"discount_percent"`
	FinalPrice      string    `json:"final_price"`
	PaymentType     string    `json:"payment_type"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func SaleToHttp(s *models.Sale) *Sale {
	return &Sale{
		ID:              s.ID,
		CarID:           s.CarID,
		ClientID:        s.ClientID,
		EmployeeID:      s.EmployeeID,
		SaleDate:        s.SaleDate.Format(TimeLayout),
		OriginPrice:     s.OriginPrice.String(),
		DiscountAmount:  s.DiscountAmount.String(),
		DiscountPercent: s.DiscountPercent.String(),
		FinalPrice:      s.FinalPrice.String(),
		PaymentType:     string(s.PaymentType),
		Status:          string(s.Status),
		Notes:           s.Notes,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

func SaleToService(s *Sale) (*models.Sale, error) {
	paymentType, ok := PaymentTypeMap[s.PaymentType]
	if !ok {
		return nil, errorx.NewError("invalid payment_type", errorx.BadRequest)
	}
	discountAmount := decimal.Zero
	if s.DiscountAmount != "" {
		var err error
		discountAmount, err = decimal.NewFromString(s.DiscountAmount)
		if err != nil {
			return nil, err
		}
	}
	return &models.Sale{
		ID:             s.ID,
		CarID:          s.CarID,
		ClientID:       s.ClientID,
		EmployeeID:     s.EmployeeID,
		DiscountAmount: discountAmount,
		PaymentType:    paymentType,
		Notes:          s.Notes,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}, nil
}
