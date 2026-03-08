package generator

import "salon/internal/models"

func (g *Generator) GenerateSupplier() (models.Supplier, error) {
	supplier, err := g.supplierStorage.PickOne()
	if err != nil {
		return models.Supplier{}, err
	}
	return *supplier, nil
}

func (g *Generator) SupplierCreated(id string, serviceSupplier *models.Supplier) error {
	return g.supplierStorage.MoveToCreated(id, serviceSupplier.ID, serviceSupplier)
}

func (g *Generator) AvailableSuppliersCount() int {
	return g.supplierStorage.AvaliableCount()
}

func (g *Generator) CreatedSuppliersCount() int {
	return g.supplierStorage.CreatedCount()
}
