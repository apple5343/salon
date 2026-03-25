package integration

import (
	"math/rand"
	"net/http"
	"salon/tests/integration/models"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	serviceModels "salon/internal/models"
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

func (s *CarSuite) TestGet() {
	s.T().Run("get", s.Get)
	s.T().Run("invalid", s.GetInvalid)
}

func (s *CarSuite) TestUpdate() {
	s.T().Run("update", s.Update)
	s.T().Run("invalid", s.UpdateInvalid)
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

func (s *CarSuite) RandomCar() *httpModels.CarInternalResponse {
	return s.cars[s.carsIDs[rand.Intn(len(s.carsIDs))]]
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

	t.Run("status sold", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Status = string(serviceModels.CarStatusSold)
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("status booked", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Status = string(serviceModels.CarStatusBooked)
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

	t.Run("existing vin", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Vin = s.RandomCar().Vin
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("empty year", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Year = 0
		_, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid year", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Year = 1072
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
	t.Run("get internal", func(t *testing.T) {
		for id := range s.cars {
			car, code, err := s.base.client.GetCar(s.base.ctx, s.managerToken.AccessToken, id)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, s.cars[id], car)
		}
	})

	t.Run("get public", func(t *testing.T) {
		for id := range s.cars {
			car, code, err := s.base.client.GetCar(s.base.ctx, "", id)
			require.NoError(t, err)
			if s.cars[id].Status != string(serviceModels.CarStatusAvailable) {
				require.Equal(t, http.StatusBadRequest, code)
			} else {
				require.Equal(t, http.StatusOK, code)
				CompareCarsPublic(t, s.cars[id], car)
			}
		}
	})
}

func (s *CarSuite) GetInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, code, err := s.base.client.GetCar(s.base.ctx, "", "invalid")
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		_, code, err := s.base.client.GetCar(s.base.ctx, "", gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *CarSuite) Update(t *testing.T) {
	for id := range s.cars {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.ID = id
		updated, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, car.ID, updated.ID)
		require.Equal(t, car.ModelID, updated.Model.ID)
		require.Equal(t, car.SupplierID, updated.Supplier.ID)
		require.Equal(t, s.cars[id].CreatedAt, updated.CreatedAt)
		require.NotEqual(t, s.cars[id].UpdatedAt, updated.UpdatedAt)
		s.cars[id] = updated
	}
}

func (s *CarSuite) UpdateInvalid(t *testing.T) {
	car := s.RandomCar()
	var car2 *httpModels.CarInternalResponse
	for id, c := range s.cars {
		if id != car.ID {
			car2 = c
			break
		}
	}
	require.NotNil(t, car2)
	t.Run("invalid id", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.ID = "invalid"
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.ID = gofakeit.UUID()
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid model id", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.ModelID = "invalid"
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing model", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.ModelID = gofakeit.UUID()
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid supplier id", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.SupplierID = "invalid"
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing supplier", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.SupplierID = gofakeit.UUID()
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty color", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Color = ""
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty interior color", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.InteriorColor = ""
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("negative mileage", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Mileage = -1
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid price", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Price = "0.0.0"
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("negative price", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Price = "-1.0"
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid status", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Status = "invalid"
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("status sold", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Status = string(serviceModels.CarStatusSold)
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("status booked", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Status = string(serviceModels.CarStatusBooked)
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty year", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Year = 0
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid year", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Year = 12345
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid vin", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Vin = car.Vin + "1"
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("existing vin", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Vin = car2.Vin
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})
}
