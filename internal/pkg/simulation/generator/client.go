package generator

import (
	"salon/internal/models"

	"github.com/brianvoe/gofakeit"
)

func (g *Generator) GenerateClient() *models.Client {
	return &models.Client{
		FullName:     gofakeit.Name(),
		Phone:        RandomPhone(),
		Email:        gofakeit.Email(),
		PasswordHash: gofakeit.Password(true, true, true, true, false, 8),
		Passport:     RandomPassport(),
		BirthDate:    gofakeit.Date(),
	}
}
