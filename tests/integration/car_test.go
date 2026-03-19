package integration

import (
	"math/rand"
	"net/http"
	"salon/tests/integration/models"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	httpModels "salon/internal/transport/http/models"
)

type CarSuite struct {
	suite.Suite
	base         *BaseTestSuite
	adminToken   *models.EmployeeToken
	managerToken *models.EmployeeToken
	suppliers    map[string]*httpModels.Supplier
	brands       map[string]*httpModels.Brand
	models       map[string]*httpModels.ModelInternalResponse
	cars         map[string]*httpModels.Car
	suppliersIDs []string
	brandsIDs    []string
	modelsIDs    []string
	carsIDs      []string
}

func (s *CarSuite) SetupSuite() {
	token, code, err := s.base.client.LoginEmployee(s.base.ctx, s.base.adminCreds.Email, s.base.adminCreds.Password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.adminToken = token

	manager := models.GenerateEmployee()
	password := manager.Password
	manager, code, err = s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, manager)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	manager, code, err = s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, manager.ID)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	token, code, err = s.base.client.LoginEmployee(s.base.ctx, manager.Email, password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.managerToken = token

	s.suppliers = make(map[string]*httpModels.Supplier)
	s.brands = make(map[string]*httpModels.Brand)
	s.models = make(map[string]*httpModels.ModelInternalResponse)
	s.cars = make(map[string]*httpModels.Car)

	for i := 0; i < 5; i++ {
		supplier, code, err := s.base.client.CreateSupplier(s.base.ctx, s.adminToken.AccessToken, models.GenerateSupplier())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		s.suppliers[supplier.ID] = supplier
		s.suppliersIDs = append(s.suppliersIDs, supplier.ID)
	}

	for i := 0; i < 5; i++ {
		brand, code, err := s.base.client.CreateBrand(s.base.ctx, s.adminToken.AccessToken, models.GenerateBrand())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		s.brands[brand.ID] = brand
		s.brandsIDs = append(s.brandsIDs, brand.ID)
	}

	for i := 0; i < 5; i++ {
		model, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, models.GenerateModel(s.RandomBrand().ID))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		s.models[model.ID] = model
		s.modelsIDs = append(s.modelsIDs, model.ID)
	}
}

func (s *CarSuite) TestCreate() {
	s.T().Run("create", s.Create)
}

func (s *CarSuite) RandomSupplier() *httpModels.Supplier {
	return s.suppliers[s.suppliersIDs[rand.Intn(len(s.suppliersIDs))]]
}

func (s *CarSuite) RandomBrand() *httpModels.Brand {
	return s.brands[s.brandsIDs[rand.Intn(len(s.brandsIDs))]]
}

func (s *CarSuite) RandomModel() *httpModels.ModelInternalResponse {
	return s.models[s.modelsIDs[rand.Intn(len(s.modelsIDs))]]
}
func (s *CarSuite) Create(t *testing.T) {
	const carsCount int = 10

	for i := 0; i < carsCount; i++ {
		car, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.cars[car.ID] = car
		s.carsIDs = append(s.carsIDs, car.ID)
	}
}
