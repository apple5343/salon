package generator

import (
	"errors"
	"salon/internal/models"
)

func (g *Generator) GenerateModel() (models.Model, error) {
	if g.modelStorage.AvaliableCount() == 0 {
		return models.Model{}, errors.New("no created models")
	}
	model, err := g.modelStorage.PickOne()
	if err != nil {
		return models.Model{}, err
	}
	return *model, nil
}

func (g *Generator) ModelCreated(id string, serviceModel *models.Model) error {
	return g.modelStorage.MoveToCreated(id, serviceModel.ID, serviceModel)
}

func (g *Generator) AvailableModelsCount() int {
	return g.modelStorage.AvaliableCount()
}

func (g *Generator) CreatedModelsCount() int {
	return g.modelStorage.CreatedCount()
}
