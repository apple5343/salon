package generator

import "salon/internal/models"

type Generator struct {
	brandStorage  *Storage[*models.Brand]
	modelStorage  *Storage[*models.Model]
	modelsByBrand map[string][]*models.Model
}

func NewGenerator() *Generator {
	g := &Generator{
		brandStorage:  NewStorage[*models.Brand](),
		modelStorage:  NewStorage[*models.Model](),
		modelsByBrand: map[string][]*models.Model{},
	}
	if err := g.LoadBrandsDataset(); err != nil {
		panic(err)
	}
	if err := g.LoadModelsDataset(); err != nil {
		panic(err)
	}
	return g
}
