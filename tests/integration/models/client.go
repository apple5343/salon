package models

import (
	httpModels "salon/internal/transport/http/models"

	"github.com/brianvoe/gofakeit"
)

func RandomDate() string {
	return gofakeit.Date().Format("02.01.2006")
}

func GenerateClient() *httpModels.Client {
	return &httpModels.Client{
		FullName:  gofakeit.Name(),
		Phone:     RandomPhone(),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, true, false, 8),
		Passport:  RandomPassport(),
		BirthDate: RandomDate(),
	}
}

func CopyClient(c *httpModels.Client) *httpModels.Client {
	return &httpModels.Client{
		ID:        c.ID,
		FullName:  c.FullName,
		Phone:     c.Phone,
		Email:     c.Email,
		Password:  c.Password,
		Passport:  c.Passport,
		BirthDate: c.BirthDate,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
