package generator

import "salon/internal/models"

type Generator struct {
	brandsAvailable map[string]*models.Brand
	brandsPending   map[string]*models.Brand
	brandsCreated   map[string]*models.Brand
	modelsByBrand   map[string][]*models.Model
}

func NewGenerator() *Generator {
	g := &Generator{
		brandsAvailable: make(map[string]*models.Brand),
		brandsPending:   make(map[string]*models.Brand),
		brandsCreated:   make(map[string]*models.Brand),
		modelsByBrand:   make(map[string][]*models.Model),
	}
	if err := g.LoadBrandsDataset(); err != nil {
		panic(err)
	}
	if err := g.LoadModelsDataset(); err != nil {
		panic(err)
	}
	return g
}
