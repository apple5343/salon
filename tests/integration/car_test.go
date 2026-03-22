package integration

import (
	"math/rand"
	"net/http"
	"salon/tests/integration/models"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	httpModels "salon/internal/transport/http/models"
)

type CarSuite struct {
	suite.Suite
	base         *BaseTestSuite
	adminToken   *models.EmployeeToken
	managerToken *models.EmployeeToken
	suppliers    map[string]*httpModels.SupplierInternalResponse
	brands       map[string]*httpModels.BrandInternalResponse
	models       map[string]*httpModels.ModelInternalResponse
	cars         map[string]*httpModels.CarInternalResponse
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

	s.suppliers = make(map[string]*httpModels.SupplierInternalResponse)
	s.brands = make(map[string]*httpModels.BrandInternalResponse)
	s.models = make(map[string]*httpModels.ModelInternalResponse)
	s.cars = make(map[string]*httpModels.CarInternalResponse)

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

	s.T().Run("create", s.Create)
}

func (s *CarSuite) TestCreate() {
	s.T().Run("create", s.Create)
	s.T().Run("invalid", s.CreateInvalid)
	s.T().Run("forbidden", s.CreateForbidden)
}

func (s *CarSuite) RandomSupplier() *httpModels.SupplierInternalResponse {
	return s.suppliers[s.suppliersIDs[rand.Intn(len(s.suppliersIDs))]]
}

func (s *CarSuite) RandomBrand() *httpModels.BrandInternalResponse {
	return s.brands[s.brandsIDs[rand.Intn(len(s.brandsIDs))]]
}

func (s *CarSuite) RandomModel() *httpModels.ModelInternalResponse {
	return s.models[s.modelsIDs[rand.Intn(len(s.modelsIDs))]]
}

func CompareCarsPublic(t *testing.T, expected, actual *httpModels.CarInternalResponse) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Vin, actual.Vin)
	require.Equal(t, expected.Year, actual.Year)
	require.Equal(t, expected.Color, actual.Color)
	require.Equal(t, expected.InteriorColor, actual.InteriorColor)
	require.Equal(t, expected.Mileage, actual.Mileage)
	require.Equal(t, expected.Price, actual.Price)
	require.Equal(t, expected.Status, expected.Status)
	CompareModelsPublic(t, expected.Model, actual.Model)
	CompareSuppliersPublic(t, expected.Supplier, actual.Supplier)
	require.Zero(t, actual.CreatedAt)
	require.Zero(t, actual.UpdatedAt)
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

func (s *CarSuite) CreateInvalid(t *testing.T) {
	t.Run("empty color", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Color = ""
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid status", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Status = "invalid"
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty interior color", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.InteriorColor = ""
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("negative mileage", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Mileage = -1
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty model id", func(t *testing.T) {
		car := models.GenerateCar("", s.RandomSupplier().ID)
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid model id", func(t *testing.T) {
		car := models.GenerateCar("invalid", s.RandomSupplier().ID)
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing model", func(t *testing.T) {
		car := models.GenerateCar(gofakeit.UUID(), s.RandomSupplier().ID)
		car.Color = ""
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid price", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Price = "0.33.33"
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("negative price", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Price = "-100.25"
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty supplier id", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, "")
		car.Color = ""
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid supplier id", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, "invalid")
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing supplier", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, gofakeit.UUID())
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty vin", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Vin = ""
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid vin", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Vin = car.Vin[:len(car.Vin)-1]
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

}

func (s *CarSuite) CreateForbidden(t *testing.T) {
	t.Run("without token", func(t *testing.T) {
		_, code, err := s.base.client.CreateCar(s.base.ctx, "", models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID))
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, code, err := s.base.client.CreateCar(s.base.ctx, gofakeit.Word(), models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID))
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("manager token", func(t *testing.T) {
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.managerToken.AccessToken, models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})
}

func (s *CarSuite) Get(t *testing.T) {

}
