package integration

import (
	"math/rand"
	"net/http"
	"salon/tests/integration/models"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/shopspring/decimal"
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

	s.T().Run("create", s.Create(serviceModels.CarStatusIncoming, 5))
}

func (s *CarSuite) TestCreate() {
	s.T().Run("create", s.Create(serviceModels.CarStatusIncoming, 5))
	s.T().Run("invalid", s.CreateInvalid)
	s.T().Run("forbidden", s.CreateForbidden)
}

func (s *CarSuite) TestGet() {
	s.T().Run("get", s.Get)
	s.T().Run("invalid", s.GetInvalid)
}

func (s *CarSuite) TestUpdate() {
	s.T().Run("update", s.Update)
	s.T().Run("get after update", s.Get)
	s.T().Run("invalid", s.UpdateInvalid)
	s.T().Run("forbidden", s.UpdateForbidden)
}

func (s *CarSuite) TestList() {
	s.T().Run("list", s.List)
	s.T().Run("invalid", s.ListInvalid)
	s.T().Run("forbidden", s.ListForbidden)
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

func CompareCarsShort(t *testing.T, expected *httpModels.CarInternalResponse, actual *httpModels.CarShort) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Vin, actual.Vin)
	require.Equal(t, expected.Year, actual.Year)
	require.Equal(t, expected.Price, actual.Price)
	require.Equal(t, expected.Status, expected.Status)
	require.Equal(t, expected.Model.Brand.Name, actual.BrandName)
	require.Equal(t, expected.Model.Name, actual.ModelName)
	require.Equal(t, expected.Supplier.Name, actual.SupplierName)
}

func (s *CarSuite) Create(status serviceModels.CarStatus, count int) func(t *testing.T) {
	return func(t *testing.T) {
		for i := 0; i < count; i++ {
			c := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
			c.Status = string(status)
			car, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, c)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			s.cars[car.ID] = car
			s.carsIDs = append(s.carsIDs, car.ID)
		}
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
	availableCount := 0
	otherCount := 0
	t.Run("get internal", func(t *testing.T) {
		for id := range s.cars {
			car, code, err := s.base.client.GetCar(s.base.ctx, s.managerToken.AccessToken, id)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, s.cars[id], car)
			if s.cars[id].Status == string(serviceModels.CarStatusAvailable) {
				availableCount++
			} else {
				otherCount++
			}
		}
	})
	if availableCount == 0 {
		t.Run("create available", s.Create(serviceModels.CarStatusAvailable, 5))
	}
	if otherCount == 0 {
		t.Run("create other", s.Create(serviceModels.CarStatusIncoming, 5))
	}
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

	var saleID string
	t.Run("update status of booked car", func(t *testing.T) {
		c := models.CarInternalToCar(car)
		c.Status = string(serviceModels.CarStatusAvailable)
		car, code, err := s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)

		client, code, err := s.base.client.RegisterClient(s.base.ctx, s.adminToken.AccessToken, models.GenerateClient())
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		sale, code, err := s.base.client.CreateSale(s.base.ctx, s.adminToken.AccessToken, models.GenerateSale(car, client.ID))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		saleID = sale.ID

		c = models.CarInternalToCar(car)
		c.Status = string(serviceModels.CarStatusPending)
		_, code, err = s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("update status of sold car", func(t *testing.T) {
		_, code, err := s.base.client.UpdateSale(s.base.ctx, s.adminToken.AccessToken, saleID, string(serviceModels.SaleStatusCompleted))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)

		car, code, err = s.base.client.GetCar(s.base.ctx, s.adminToken.AccessToken, car.ID)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, string(serviceModels.CarStatusSold), car.Status)
		s.cars[car.ID] = car

		c := models.CarInternalToCar(car)
		c.Status = string(serviceModels.CarStatusPending)
		_, code, err = s.base.client.UpdateCar(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	//TODO обноаление машины с статусом sold или booked
}

func (s *CarSuite) UpdateForbidden(t *testing.T) {
	car := models.CarInternalToCar(s.RandomCar())
	t.Run("without token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateCar(s.base.ctx, "", car)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with invalid token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateCar(s.base.ctx, gofakeit.Word(), car)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with manager token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateCar(s.base.ctx, s.managerToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})
}

func (s *CarSuite) CheckList(t *testing.T, token string, filter *serviceModels.CarFilters, expected []string) (map[string]*httpModels.CarShort, []string) {
	all := make(map[string]*httpModels.CarShort)
	sorted := []string{}
	for true {
		cars, code, err := s.base.client.GetCars(s.base.ctx, token, filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.LessOrEqual(t, len(cars), *filter.Limit)
		for _, c := range cars {
			id := c.ID
			if _, ok := s.cars[id]; !ok {
				t.Errorf("invalid car id: %s", id)
				continue
			}
			if _, ok := all[id]; ok {
				t.Errorf("duplicate car id: %s", id)
				continue
			}
			all[id] = c
			sorted = append(sorted, id)
		}
		if len(cars) < *filter.Limit {
			break
		}
		*filter.Offset += *filter.Limit
	}
	for _, id := range expected {
		if _, ok := all[id]; !ok {
			t.Errorf("not checked car id: %s", id)
		}
	}
	require.Equal(t, len(expected), len(all))
	return all, sorted
}

func (s *CarSuite) List(t *testing.T) {
	statuses := []serviceModels.CarStatus{serviceModels.CarStatusAvailable, serviceModels.CarStatusIncoming,
		serviceModels.CarStatusPending, serviceModels.CarStatusArchived}
	for _, status := range statuses {
		t.Run("create "+string(status), s.Create(status, 5))
	}

	t.Run("all cars & order by created at as asc", func(t *testing.T) {
		// сортировка по created_at asc должна применяться по умолчанию
		limit := 6
		offset := 0
		filter := serviceModels.CarFilters{
			BaseList: serviceModels.BaseList{
				Limit:  &limit,
				Offset: &offset,
			},
		}
		expected := make([]string, 0, len(s.cars))
		for id := range s.cars {
			expected = append(expected, id)
		}
		all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
		for i := 1; i < len(sorted); i++ {
			if s.cars[sorted[i-1]].CreatedAt.After(s.cars[sorted[i]].CreatedAt) {
				t.Errorf("invalid order by created at asc: %v > %v", s.cars[sorted[i-1]].CreatedAt, s.cars[sorted[i]].CreatedAt)
			}
		}
		for id, c := range all {
			CompareCarsShort(t, s.cars[id], c)
		}
	})

	t.Run("all available cars & order by price desc", func(t *testing.T) {
		limit := 8
		offset := 0
		orderDirection := serviceModels.OrderDirectionDESC
		orderBy := serviceModels.CarOrderByPrice
		// для неавторизованных пользователей доступны только активные машины (available status)
		filter := serviceModels.CarFilters{
			BaseList: serviceModels.BaseList{
				Limit:          &limit,
				Offset:         &offset,
				OrderDirection: &orderDirection,
			},
			OrderBy: &orderBy,
		}
		expected := []string{}
		for id := range s.cars {
			if s.cars[id].Status == string(serviceModels.CarStatusAvailable) {
				expected = append(expected, id)
			}
		}
		all, sorted := s.CheckList(t, "", &filter, expected)
		for i := 1; i < len(sorted); i++ {
			price1, err := decimal.NewFromString(s.cars[sorted[i-1]].Price)
			require.NoError(t, err)
			price2, err := decimal.NewFromString(s.cars[sorted[i]].Price)
			require.NoError(t, err)
			if price1.LessThan(price2) {
				t.Errorf("invalid order by price desc: %v < %v", s.cars[sorted[i-1]].Price, s.cars[sorted[i]].Price)
			}
		}
		for id, c := range all {
			CompareCarsShort(t, s.cars[id], c)
		}
	})

	type listFunc func() (string, *serviceModels.CarFilters, []string)
	type checked func(t *testing.T, all map[string]*httpModels.CarShort, sorted []string)
	tests := []struct {
		name  string
		list  listFunc
		check checked
	}{
		{
			name: "restricted price range & order by updated at desc",
			list: func() (string, *serviceModels.CarFilters, []string) {
				limit := 12
				offset := 0
				orderDirection := serviceModels.OrderDirectionDESC
				orderBy := serviceModels.CarOrderByUpdatedAt
				minPrice := gofakeit.Number(100, 1000)
				maxPrice := minPrice + gofakeit.Number(50, 200)
				minPriceDecimal, err := decimal.NewFromString(strconv.Itoa(minPrice))
				require.NoError(t, err)
				maxPriceDecimal, err := decimal.NewFromString(strconv.Itoa(maxPrice))
				require.NoError(t, err)
				filter := serviceModels.CarFilters{
					BaseList: serviceModels.BaseList{
						Limit:          &limit,
						Offset:         &offset,
						OrderDirection: &orderDirection,
					},
					MinPrice: &minPriceDecimal,
					MaxPrice: &maxPriceDecimal,
					OrderBy:  &orderBy,
				}
				expected := []string{}
				for id := range s.cars {
					price, err := decimal.NewFromString(s.cars[id].Price)
					require.NoError(t, err)
					if price.GreaterThanOrEqual(minPriceDecimal) && price.LessThanOrEqual(maxPriceDecimal) {
						expected = append(expected, id)
					}
				}
				return s.adminToken.AccessToken, &filter, expected
			},
			check: func(t *testing.T, all map[string]*httpModels.CarShort, sorted []string) {
				for i := 1; i < len(sorted); i++ {
					if s.cars[sorted[i-1]].UpdatedAt.Before(s.cars[sorted[i]].UpdatedAt) {
						t.Errorf("invalid order by updated at desc: %v < %v", s.cars[sorted[i-1]].UpdatedAt, s.cars[sorted[i]].UpdatedAt)
					}
				}
				for id, c := range all {
					CompareCarsShort(t, s.cars[id], c)
				}
			},
		},
		{
			name: "restricted mileage range & order by mileage asc",
			list: func() (string, *serviceModels.CarFilters, []string) {
				limit := 12
				offset := 0
				orderBy := serviceModels.CarOrderByMileage

				min := 20
				max := 60

				filter := serviceModels.CarFilters{
					BaseList: serviceModels.BaseList{
						Limit:  &limit,
						Offset: &offset,
					},
					MinMileage: &min,
					MaxMileage: &max,
					OrderBy:    &orderBy,
				}
				expected := []string{}
				for id, car := range s.cars {
					if car.Mileage >= min && car.Mileage <= max {
						expected = append(expected, id)
					}
				}
				return s.adminToken.AccessToken, &filter, expected
			},
			check: func(t *testing.T, all map[string]*httpModels.CarShort, sorted []string) {
				for i := 1; i < len(sorted); i++ {
					if s.cars[sorted[i-1]].Mileage > s.cars[sorted[i]].Mileage {
						t.Errorf("invalid order by mileage asc: %v > %v", s.cars[sorted[i-1]].Mileage, s.cars[sorted[i]].Mileage)
					}
				}
				for id, c := range all {
					CompareCarsShort(t, s.cars[id], c)
				}
			},
		},
		{
			name: "restricted year range & order by year desc",
			list: func() (string, *serviceModels.CarFilters, []string) {
				limit := 12
				offset := 0
				orderDirection := serviceModels.OrderDirectionDESC
				orderBy := serviceModels.CarOrderByYear
				min := 2005
				max := 2015

				filter := serviceModels.CarFilters{
					BaseList: serviceModels.BaseList{
						Limit:          &limit,
						Offset:         &offset,
						OrderDirection: &orderDirection,
					},
					MinYear: &min,
					MaxYear: &max,
					OrderBy: &orderBy,
				}
				expected := []string{}
				for id, car := range s.cars {
					if car.Year >= min && car.Year <= max {
						expected = append(expected, id)
					}
				}
				return s.adminToken.AccessToken, &filter, expected
			},
			check: func(t *testing.T, all map[string]*httpModels.CarShort, sorted []string) {
				for i := 1; i < len(sorted); i++ {
					if s.cars[sorted[i-1]].Year < s.cars[sorted[i]].Year {
						t.Errorf("invalid order by year desc: %v < %v", s.cars[sorted[i-1]].Year, s.cars[sorted[i]].Year)
					}
				}
				for id, c := range all {
					CompareCarsShort(t, s.cars[id], c)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, filter, expected := tt.list()
			all, sorted := s.CheckList(t, token, filter, expected)
			tt.check(t, all, sorted)
		})
	}

	for _, status := range statuses {
		t.Run("filter by status "+string(status)+" & order by created at asc", func(t *testing.T) {
			limit := 10
			offset := 0
			filter := serviceModels.CarFilters{
				BaseList: serviceModels.BaseList{
					Limit:  &limit,
					Offset: &offset,
				},
				Status: &status,
			}
			expected := []string{}
			for id := range s.cars {
				if s.cars[id].Status == string(status) {
					expected = append(expected, id)
				}
			}
			all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
			for i := 1; i < len(sorted); i++ {
				if s.cars[sorted[i-1]].CreatedAt.After(s.cars[sorted[i]].CreatedAt) {
					t.Errorf("invalid order by created at asc: %v > %v", s.cars[sorted[i-1]].CreatedAt, s.cars[sorted[i]].CreatedAt)
				}
			}
			for id, c := range all {
				CompareCarsShort(t, s.cars[id], c)
			}
		})
	}

	t.Run("filter by supplier_id & order by year desc", func(t *testing.T) {
		for supplierID := range s.suppliers {
			limit := 10
			offset := 0
			orderBy := serviceModels.CarOrderByYear
			orderDirection := serviceModels.OrderDirectionDESC
			filter := serviceModels.CarFilters{
				BaseList: serviceModels.BaseList{
					Limit:          &limit,
					Offset:         &offset,
					OrderDirection: &orderDirection,
				},
				SupplierID: &supplierID,
				OrderBy:    &orderBy,
			}
			expected := []string{}
			for id := range s.cars {
				if s.cars[id].Supplier.ID == supplierID {
					expected = append(expected, id)
				}
			}
			if len(expected) == 0 {
				continue
			}
			all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
			for i := 1; i < len(sorted); i++ {
				if s.cars[sorted[i-1]].Year < s.cars[sorted[i]].Year {
					t.Errorf("invalid order by year desc: %v < %v", s.cars[sorted[i-1]].Year, s.cars[sorted[i]].Year)
				}
			}
			for id, c := range all {
				require.Equal(t, s.suppliers[supplierID].Name, c.SupplierName)
				CompareCarsShort(t, s.cars[id], c)
			}
		}
	})

	t.Run("filter by model_id & order by mileage asc", func(t *testing.T) {
		for modelID := range s.models {
			limit := 10
			offset := 0
			orderBy := serviceModels.CarOrderByMileage
			filter := serviceModels.CarFilters{
				BaseList: serviceModels.BaseList{
					Limit:  &limit,
					Offset: &offset,
				},
				ModelID: &modelID,
				OrderBy: &orderBy,
			}
			expected := []string{}
			for id := range s.cars {
				if s.cars[id].Model.ID == modelID {
					expected = append(expected, id)
				}
			}
			if len(expected) == 0 {
				continue
			}
			all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
			for i := 1; i < len(sorted); i++ {
				if s.cars[sorted[i-1]].Mileage > s.cars[sorted[i]].Mileage {
					t.Errorf("invalid order by mileage asc: %v > %v", s.cars[sorted[i-1]].Mileage, s.cars[sorted[i]].Mileage)
				}
			}
			for id, c := range all {
				require.Equal(t, s.models[modelID].Name, c.ModelName)
				CompareCarsShort(t, s.cars[id], c)
			}
		}
	})

	t.Run("filter by brand_id & order by price asc", func(t *testing.T) {
		for brandID := range s.brands {
			limit := 10
			offset := 0
			orderBy := serviceModels.CarOrderByPrice
			filter := serviceModels.CarFilters{
				BaseList: serviceModels.BaseList{
					Limit:  &limit,
					Offset: &offset,
				},
				BrandID: &brandID,
				OrderBy: &orderBy,
			}
			expected := []string{}
			for id := range s.cars {
				if s.cars[id].Model.Brand.ID == brandID {
					expected = append(expected, id)
				}
			}
			if len(expected) == 0 {
				continue
			}
			all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
			for i := 1; i < len(sorted); i++ {
				price1, err := decimal.NewFromString(s.cars[sorted[i-1]].Price)
				require.NoError(t, err)
				price2, err := decimal.NewFromString(s.cars[sorted[i]].Price)
				require.NoError(t, err)
				if price1.GreaterThan(price2) {
					t.Errorf("invalid order by price asc: %v > %v", s.cars[sorted[i-1]].Price, s.cars[sorted[i]].Price)
				}
			}
			for id, c := range all {
				require.Equal(t, s.brands[brandID].Name, c.BrandName)
				CompareCarsShort(t, s.cars[id], c)
			}
		}
	})

	t.Run("filter by brand_id and year range & order by price desc", func(t *testing.T) {
		for brandID := range s.brands {
			limit := 10
			offset := 0
			orderBy := serviceModels.CarOrderByPrice
			orderDirection := serviceModels.OrderDirectionDESC
			minYear := 2010
			maxYear := 2020
			filter := serviceModels.CarFilters{
				BaseList: serviceModels.BaseList{
					Limit:          &limit,
					Offset:         &offset,
					OrderDirection: &orderDirection,
				},
				BrandID: &brandID,
				MinYear: &minYear,
				MaxYear: &maxYear,
				OrderBy: &orderBy,
			}
			expected := []string{}
			for id := range s.cars {
				if s.cars[id].Model.Brand.ID == brandID && s.cars[id].Year >= minYear && s.cars[id].Year <= maxYear {
					expected = append(expected, id)
				}
			}
			if len(expected) == 0 {
				continue
			}
			all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
			for i := 1; i < len(sorted); i++ {
				price1, err := decimal.NewFromString(s.cars[sorted[i-1]].Price)
				require.NoError(t, err)
				price2, err := decimal.NewFromString(s.cars[sorted[i]].Price)
				require.NoError(t, err)
				if price1.LessThan(price2) {
					t.Errorf("invalid order by price desc: %v < %v", s.cars[sorted[i-1]].Price, s.cars[sorted[i]].Price)
				}
			}
			for id, c := range all {
				require.Equal(t, s.brands[brandID].Name, c.BrandName)
				require.GreaterOrEqual(t, c.Year, minYear)
				require.LessOrEqual(t, c.Year, maxYear)
				CompareCarsShort(t, s.cars[id], c)
			}
		}
	})

	t.Run("filter by supplier_id and mileage range & order by year asc", func(t *testing.T) {
		for supplierID := range s.suppliers {
			limit := 10
			offset := 0
			orderBy := serviceModels.CarOrderByYear
			minMileage := 10
			maxMileage := 80
			filter := serviceModels.CarFilters{
				BaseList: serviceModels.BaseList{
					Limit:  &limit,
					Offset: &offset,
				},
				SupplierID: &supplierID,
				MinMileage: &minMileage,
				MaxMileage: &maxMileage,
				OrderBy:    &orderBy,
			}
			expected := []string{}
			for id := range s.cars {
				if s.cars[id].Supplier.ID == supplierID && s.cars[id].Mileage >= minMileage && s.cars[id].Mileage <= maxMileage {
					expected = append(expected, id)
				}
			}
			if len(expected) == 0 {
				continue
			}
			all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
			for i := 1; i < len(sorted); i++ {
				if s.cars[sorted[i-1]].Year > s.cars[sorted[i]].Year {
					t.Errorf("invalid order by year asc: %v > %v", s.cars[sorted[i-1]].Year, s.cars[sorted[i]].Year)
				}
			}
			for id, c := range all {
				require.Equal(t, s.suppliers[supplierID].Name, c.SupplierName)
				require.GreaterOrEqual(t, s.cars[id].Mileage, minMileage)
				require.LessOrEqual(t, s.cars[id].Mileage, maxMileage)
				CompareCarsShort(t, s.cars[id], c)
			}
		}
	})

	t.Run("filter by all combinations of color and interior color & order by price asc", func(t *testing.T) {
		colors := models.CarColors
		limit := 50
		orderBy := serviceModels.CarOrderByPrice
		orderDirection := serviceModels.OrderDirectionASC

		carsByColor := make(map[string]map[string][]string)
		for _, color := range colors {
			carsByColor[color] = make(map[string][]string)
			for _, intCol := range colors {
				carsByColor[color][intCol] = []string{}
			}
		}
		for id, car := range s.cars {
			carsByColor[car.Color][car.InteriorColor] = append(carsByColor[car.Color][car.InteriorColor], id)
		}

		for _, color := range colors {
			for _, intCol := range colors {
				expected := carsByColor[color][intCol]
				if len(expected) == 0 {
					continue
				}
				offset := 0
				filter := serviceModels.CarFilters{
					BaseList: serviceModels.BaseList{
						Limit:          &limit,
						Offset:         &offset,
						OrderDirection: &orderDirection,
					},
					OrderBy:       &orderBy,
					Color:         &color,
					InteriorColor: &intCol,
				}
				all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
				for i := 1; i < len(sorted); i++ {
					pricePrev, err := decimal.NewFromString(s.cars[sorted[i-1]].Price)
					require.NoError(t, err)
					priceCur, err := decimal.NewFromString(s.cars[sorted[i]].Price)
					require.NoError(t, err)
					if pricePrev.GreaterThan(priceCur) {
						t.Errorf("Order by price asc failed for color=%s, interior_color=%s: %v > %v", color, intCol, s.cars[sorted[i-1]].Price, s.cars[sorted[i]].Price)
					}
				}
				for id, c := range all {
					CompareCarsShort(t, s.cars[id], c)
				}
			}
		}
	})
}

func (s *CarSuite) ListInvalid(t *testing.T) {
	t.Run("with invalid order_by", func(t *testing.T) {
		var orderBy serviceModels.CarOrderBy = "invalid"
		filter := serviceModels.CarFilters{OrderBy: &orderBy}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with invalid order_direction", func(t *testing.T) {
		var orderDirection serviceModels.OrderDirection = "invalid"
		filter := serviceModels.CarFilters{
			BaseList: serviceModels.BaseList{
				OrderDirection: &orderDirection,
			},
		}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with invalid status", func(t *testing.T) {
		var status serviceModels.CarStatus = "invalid"
		filter := serviceModels.CarFilters{Status: &status}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})


	t.Run("with negative limit", func(t *testing.T) {
		limit := -10
		filter := serviceModels.CarFilters{
			BaseList: serviceModels.BaseList{
				Limit: &limit,
			},
		}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with negative offset", func(t *testing.T) {
		offset := -10
		filter := serviceModels.CarFilters{
			BaseList: serviceModels.BaseList{
				Offset: &offset,
			},
		}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})


	t.Run("with negative min_mileage", func(t *testing.T) {
		minMileage := -10
		filter := serviceModels.CarFilters{MinMileage: &minMileage}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with negative max_mileage", func(t *testing.T) {
		maxMileage := -10
		filter := serviceModels.CarFilters{MaxMileage: &maxMileage}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with min_mileage greater than max_mileage", func(t *testing.T) {
		minMileage := 100000
		maxMileage := 50000
		filter := serviceModels.CarFilters{MinMileage: &minMileage, MaxMileage: &maxMileage}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with invalid min_year", func(t *testing.T) {
		minYear := 1800
		filter := serviceModels.CarFilters{MinYear: &minYear}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with min_year greater than max_year", func(t *testing.T) {
		minYear := 2020
		maxYear := 2010
		filter := serviceModels.CarFilters{MinYear: &minYear, MaxYear: &maxYear}
		_, code, err := s.base.client.GetCars(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *CarSuite) ListForbidden(t *testing.T) {
	t.Skip("Cars list is public endpoint")
}
