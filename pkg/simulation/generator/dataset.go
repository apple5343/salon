package generator

import (
	"encoding/json"
	"os"
	repo "salon/internal/repository/models"
)

func (g *Generator) LoadBrandsDataset() error {
	data, err := os.ReadFile("datasets/brands.json") //TODO parse path from env
	if err != nil {
		return err
	}
	var brands []*repo.Brand
	err = json.Unmarshal(data, &brands)
	if err != nil {
		return err
	}
	for _, brand := range brands {
		g.brandsAvailable[brand.ID] = repo.BrandToService(brand)
	}
	return nil
}

func (g *Generator) LoadModelsDataset() error {
	data, err := os.ReadFile("datasets/models.json") //TODO parse path from env
	if err != nil {
		return err
	}
	var models []*repo.Model
	err = json.Unmarshal(data, &models)
	if err != nil {
		return err
	}
	for _, model := range models {
		g.modelsByBrand[model.BrandID] = append(g.modelsByBrand[model.BrandID], repo.ModelToService(model))
	}
	return nil
}
