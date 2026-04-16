package generator

import (
	"salon/internal/models"
)

func (g *Generator) GenerateBrand() (models.Brand, error) {
	brand, err := g.brandStorage.PickOne()
	if err != nil {
		return models.Brand{}, err
	}
	return *brand, nil
}

func (g *Generator) BrandCreated(id string, serviceBrand *models.Brand) error {
	if err := g.brandStorage.MoveToCreated(id, serviceBrand.ID, serviceBrand); err != nil {
		return err
	}
	for _, model := range g.modelsByBrand[id] {
		model.BrandID = serviceBrand.ID
		g.modelStorage.MoveToAvailable(model.ID, model.ID, model)
	}
	return nil
}

func (g *Generator) AvailableBrandsCount() int {
	return g.brandStorage.AvaliableCount()
}

func (g *Generator) CreatedBrandsCount() int {
	return g.brandStorage.CreatedCount()
}
