package generator

import (
	"errors"
	"math/rand/v2"
	"salon/internal/models"
)

func (g *Generator) GenerateBrand() (models.Brand, error) {
	if len(g.brandsAvailable) == 0 {
		return models.Brand{}, errors.New("no available brands")
	}
	keys := make([]string, 0, len(g.brandsAvailable))
	for k := range g.brandsAvailable {
		keys = append(keys, k)
	}
	id := keys[rand.IntN(len(keys))]
	g.brandsPending[id] = g.brandsAvailable[id]
	delete(g.brandsAvailable, id)
	return *g.brandsPending[id], nil
}

func (g *Generator) BrandCreated(id string, serviceID string) error {
	if _, ok := g.brandsPending[id]; !ok {
		return errors.New("brand not pending")
	}
	brand := g.brandsPending[id]
	brand.ID = serviceID
	g.brandsCreated[id] = brand
	delete(g.brandsPending, id)
	return nil
}

func (g *Generator) AvailableBrandsCount() int {
	return len(g.brandsAvailable)
}
