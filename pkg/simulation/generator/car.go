package generator

import (
	"salon/internal/models"

	"github.com/brianvoe/gofakeit"
)

func (g *Generator) GenerateCar() models.Car {
	return models.Car{
		ID: gofakeit.UUID(),
	}
}
