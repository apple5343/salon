package generator

import "salon/internal/models"

type CarState struct {
	modelServiceID    string
	supplierServiceID string
}


// Генерирует заданные в датасетах данные для симуляции с учетом их связей.Также хранит состояния сгенерированных сущностей для имитации Foreign Keys
type Generator struct {
	brandStorage    *Storage[*models.Brand]
	modelStorage    *Storage[*models.Model]
	supplierStorage *Storage[*models.Supplier]
	carStorage      *Storage[*models.Car]

	modelsByBrand  map[string][]*models.Model
	carsByModel    map[string][]string
	carsBySupplier map[string][]string

	carStates map[string]CarState
}

func NewGenerator() *Generator {
	g := &Generator{
		brandStorage:    NewStorage[*models.Brand](),
		modelStorage:    NewStorage[*models.Model](),
		supplierStorage: NewStorage[*models.Supplier](),
		carStorage:      NewStorage[*models.Car](),
		modelsByBrand:   map[string][]*models.Model{},
		carsByModel:     map[string][]string{},
		carsBySupplier:  map[string][]string{},
		carStates:       map[string]CarState{},
	}
	if err := g.LoadBrandsDataset(); err != nil {
		panic(err)
	}
	if err := g.LoadModelsDataset(); err != nil {
		panic(err)
	}
	if err := g.LoadSuppliersDataset(); err != nil {
		panic(err)
	}
	if err := g.LoadCarsDataset(); err != nil {
		panic(err)
	}
	return g
}
