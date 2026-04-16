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
		g.brandStorage.AddAvailable(brand.ID, repo.BrandToService(brand))
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
		m := repo.ModelToService(model)
		g.modelStorage.AddPending(model.ID, m)
		g.modelsByBrand[model.BrandID] = append(g.modelsByBrand[model.BrandID], m)
	}
	return nil
}

func (g *Generator) LoadSuppliersDataset() error {
	data, err := os.ReadFile("datasets/suppliers.json") //TODO parse path from env
	if err != nil {
		return err
	}
	var suppliers []*repo.Supplier
	err = json.Unmarshal(data, &suppliers)
	if err != nil {
		return err
	}
	for _, supplier := range suppliers {
		g.supplierStorage.AddAvailable(supplier.ID, repo.SupplierToService(supplier))
	}
	return nil
}

func (g *Generator) LoadCarsDataset() error {
	data, err := os.ReadFile("datasets/cars.json") //TODO parse path from env
	if err != nil {
		return err
	}
	var cars []*repo.Car
	err = json.Unmarshal(data, &cars)
	if err != nil {
		return err
	}
	for _, car := range cars {
		g.carsByModel[car.ModelID] = append(g.carsByModel[car.ModelID], car.ID)
		g.carsBySupplier[car.SupplierID] = append(g.carsBySupplier[car.SupplierID], car.ID)
		g.carStates[car.ID] = CarState{}
		g.carStorage.AddPending(car.ID, repo.CarToService(car))
	}
	return nil
}
