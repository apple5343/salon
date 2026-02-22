package generator

import (
	"errors"
	"math/rand/v2"
	"salon/internal/models"
)

func (g *Generator) GenerateModel() (models.Model, error) {
	if len(g.brandsCreated) == 0 {
		return models.Model{}, errors.New("no created brands")
	}
	keys := make([]string, 0, len(g.brandsCreated))
	for k := range g.brandsCreated {
		keys = append(keys, k)
	}
	id := keys[rand.IntN(len(keys))]
	m, ok := g.modelsByBrand[id]
	if !ok {
		return models.Model{}, errors.New("no available models for brand")
	}
	model := *m[rand.IntN(len(m))]
	model.BrandID = g.brandsCreated[id].ID
	return model, nil
}
