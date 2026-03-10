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
	for _, carID := range g.carsByModel[id] {
		state, ok := g.carStates[carID]
		if !ok {
			continue
		}
		state.modelServiceID = serviceModel.ID
		g.carStates[carID] = state
		g.checkCarState(carID)
	}
	return g.modelStorage.MoveToCreated(id, serviceModel.ID, serviceModel)
}

func (g *Generator) AvailableModelsCount() int {
	return g.modelStorage.AvaliableCount()
}

func (g *Generator) PendingModelsCount() int {
	return g.modelStorage.PendingCount()
}

func (g *Generator) CreatedModelsCount() int {
	return g.modelStorage.CreatedCount()
}
