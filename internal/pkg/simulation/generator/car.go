package generator

import (
	"errors"
	"fmt"
	"math/rand"
	"salon/internal/models"
)

func generateNewVIN(baseVIN string) string {
	prefix := baseVIN[:11]
	newSuffix := fmt.Sprintf("%06d", rand.Intn(1000000))
	return prefix + newSuffix
}

func (g *Generator) checkCarState(carID string) error {
	state, ok := g.carStates[carID]
	if !ok {
		return errors.New("car not found")
	}
	if state.modelServiceID != "" && state.supplierServiceID != "" {
		car, ok := g.carStorage.GetPending(carID)
		if !ok {
			return errors.New("car not found in pending storage")
		}
		car.SupplierID = state.supplierServiceID
		car.ModelID = state.modelServiceID
		g.carStorage.MoveToAvailable(carID, car.ID, car)
	}
	return nil
}

func (g *Generator) GenerateCar() (models.Car, error) {
	car, err := g.carStorage.GetOneAvailable()
	if err != nil {
		return models.Car{}, err
	}
	car.Status = models.CarStatusIncoming
	car.Vin = generateNewVIN(car.Vin)
	return *car, nil
}

func (g *Generator) CarCreated(serviceCar *models.Car) {
	g.carStorage.AddCreated(serviceCar.ID, serviceCar)
}

func (g *Generator) GetCreatedCar(id string) (*models.Car, bool) {
	return g.carStorage.GetCreated(id)
}

func (g *Generator) CarUpdated(car *models.Car) {
	g.carStorage.AddCreated(car.ID, car)
}

func (g *Generator) AvailableCarsCount() int {
	return g.carStorage.AvaliableCount()
}

func (g *Generator) CreatedCarsCount() int {
	return g.carStorage.CreatedCount()
}
