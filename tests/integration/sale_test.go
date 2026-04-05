package integration

import (
	"net/http"
	"salon/tests/integration/models"
	"testing"

	serviceModels "salon/internal/models"
	httpModels "salon/internal/transport/http/models"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SaleSuite struct {
	suite.Suite
	base            *BaseTestSuite
	adminToken      *models.EmployeeToken
	clients         map[string]*httpModels.Client
	employees       map[string]*httpModels.Employee
	employeesTokens map[string]*models.EmployeeToken
	brands          map[string]*httpModels.BrandInternalResponse
	models          map[string]*httpModels.ModelInternalResponse
	cars            map[string]*httpModels.CarInternalResponse
	suppliers       map[string]*httpModels.SupplierInternalResponse
	sales           map[string]*httpModels.Sale
	salesIDs        []string
	brandsIDs       []string
	modelsIDs       []string
	carsIDs         []string
	suppliersIDs    []string
	employeesIDs    []string
	clientsIDs      []string
}

func (s *SaleSuite) SetupSuite() {
	token, code, err := s.base.client.LoginEmployee(s.base.ctx, s.base.adminCreds.Email, s.base.adminCreds.Password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.adminToken = token

	s.clients = make(map[string]*httpModels.Client)
	s.employees = make(map[string]*httpModels.Employee)
	s.employeesTokens = make(map[string]*models.EmployeeToken)
	s.brands = make(map[string]*httpModels.BrandInternalResponse)
	s.models = make(map[string]*httpModels.ModelInternalResponse)
	s.cars = make(map[string]*httpModels.CarInternalResponse)
	s.suppliers = make(map[string]*httpModels.SupplierInternalResponse)
	s.sales = make(map[string]*httpModels.Sale)

	for i := 0; i < 5; i++ {
		manager := models.GenerateEmployee()
		password := manager.Password
		manager, code, err = s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, manager)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		manager, code, err = s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, manager.ID)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		employeeToken, code, err := s.base.client.LoginEmployee(s.base.ctx, manager.Email, password)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		s.employees[manager.ID] = manager
		s.employeesIDs = append(s.employeesIDs, manager.ID)
		s.employeesTokens[manager.ID] = employeeToken
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

	for i := 0; i < 5; i++ {
		supplier, code, err := s.base.client.CreateSupplier(s.base.ctx, s.adminToken.AccessToken, models.GenerateSupplier())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		s.suppliers[supplier.ID] = supplier
		s.suppliersIDs = append(s.suppliersIDs, supplier.ID)
	}

	for i := 0; i < 5; i++ {
		client, code, err := s.base.client.RegisterClient(s.base.ctx, s.adminToken.AccessToken, models.GenerateClient())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		s.clients[client.ID] = client
		s.clientsIDs = append(s.clientsIDs, client.ID)
	}
}

func (s *SaleSuite) RandomBrand() *httpModels.BrandInternalResponse {
	return s.brands[s.brandsIDs[gofakeit.Number(0, len(s.brandsIDs)-1)]]
}

func (s *SaleSuite) RandomModel() *httpModels.ModelInternalResponse {
	return s.models[s.modelsIDs[gofakeit.Number(0, len(s.modelsIDs)-1)]]
}

func (s *SaleSuite) RandomSupplier() *httpModels.SupplierInternalResponse {
	return s.suppliers[s.suppliersIDs[gofakeit.Number(0, len(s.suppliersIDs)-1)]]
}

func (s *SaleSuite) RandomEmployee() *httpModels.Employee {
	return s.employees[s.employeesIDs[gofakeit.Number(0, len(s.employeesIDs)-1)]]
}

func (s *SaleSuite) RandomCar() *httpModels.CarInternalResponse {
	return s.cars[s.carsIDs[gofakeit.Number(0, len(s.carsIDs)-1)]]
}

func (s *SaleSuite) RandomClient() *httpModels.Client {
	return s.clients[s.clientsIDs[gofakeit.Number(0, len(s.clientsIDs)-1)]]
}

func (s *SaleSuite) TestCreate() {
	s.T().Run("complete", s.CompleteSaleFlow)
	s.T().Run("cancel", s.CancelSaleFlow)
	s.T().Run("get", s.Get)
	s.T().Run("invalid", s.CreateInvalid)
	s.T().Run("forbidden", s.CreateForbidden)
}

func (s *SaleSuite) TestGet() {
	s.T().Run("get", s.Get)
	s.T().Run("invalid", s.GetInvalid)
	s.T().Run("forbidden", s.GetForbidden)
}

func (s *SaleSuite) TestUpdate() {
	s.T().Run("invalid", s.UpdateInvalid)
	s.T().Run("forbidden", s.UpdateForbidden)
}

func (s *SaleSuite) CreateSales(sales []*httpModels.Sale) func(t *testing.T) {
	return func(t *testing.T) {
		for i := range sales {
			car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
			car.Status = string(serviceModels.CarStatusAvailable)
			carInternal, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.cars[carInternal.ID] = carInternal
			s.carsIDs = append(s.carsIDs, carInternal.ID)

			employee := s.RandomEmployee()
			token := s.employeesTokens[employee.ID].AccessToken

			sale, code, err := s.base.client.CreateSale(s.base.ctx, token, models.GenerateSale(carInternal, s.RandomClient().ID))
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.SaleStatusPending), sale.Status)
			s.sales[sale.ID] = sale
			s.salesIDs = append(s.salesIDs, sale.ID)
			sales[i] = sale
		}
	}
}

func (s *SaleSuite) CompleteSaleFlow(t *testing.T) {
	const salesCount int = 10
	sales := make([]*httpModels.Sale, salesCount)

	t.Run("create", s.CreateSales(sales))

	t.Run("check cars statuses after creation", func(t *testing.T) {
		for _, sale := range sales {
			// Обновить статус продажи может как и работник, создавший продажу, так и админ
			car, code, err := s.base.client.GetCar(s.base.ctx, s.adminToken.AccessToken, sale.CarID)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.CarStatusBooked), car.Status)
		}
	})

	t.Run("complete", func(t *testing.T) {
		for _, sale := range sales {
			sale, code, err := s.base.client.UpdateSale(s.base.ctx, s.adminToken.AccessToken, sale.ID, string(serviceModels.SaleStatusCompleted))
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.SaleStatusCompleted), sale.Status)
			s.sales[sale.ID] = sale
		}
	})

	t.Run("check cars statuses after completion", func(t *testing.T) {
		for _, sale := range sales {
			car, code, err := s.base.client.GetCar(s.base.ctx, s.adminToken.AccessToken, sale.CarID)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.CarStatusSold), car.Status)
		}
	})

	t.Run("check sale statuses after completion", func(t *testing.T) {
		for _, sale := range sales {
			sale, code, err := s.base.client.GetSale(s.base.ctx, s.adminToken.AccessToken, sale.ID)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.SaleStatusCompleted), sale.Status)
		}
	})
}

func (s *SaleSuite) CancelSaleFlow(t *testing.T) {
	const salesCount int = 10
	sales := make([]*httpModels.Sale, salesCount)

	t.Run("create", s.CreateSales(sales))

	t.Run("check cars statuses after creation", func(t *testing.T) {
		for _, sale := range sales {
			car, code, err := s.base.client.GetCar(s.base.ctx, s.adminToken.AccessToken, sale.CarID)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.CarStatusBooked), car.Status)
		}
	})

	t.Run("cancel", func(t *testing.T) {
		for _, sale := range sales {
			sale, code, err := s.base.client.UpdateSale(s.base.ctx, s.employeesTokens[sale.EmployeeID].AccessToken, sale.ID, string(serviceModels.SaleStatusCanceled))
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.SaleStatusCanceled), sale.Status)
			s.sales[sale.ID] = sale
		}
	})

	t.Run("check cars statuses after cancellation", func(t *testing.T) {
		for _, sale := range sales {
			car, code, err := s.base.client.GetCar(s.base.ctx, s.adminToken.AccessToken, sale.CarID)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.CarStatusAvailable), car.Status)
		}
	})

	t.Run("check sales statuses after cancellation", func(t *testing.T) {
		for _, sale := range sales {
			sale, code, err := s.base.client.GetSale(s.base.ctx, s.adminToken.AccessToken, sale.ID)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(string(serviceModels.SaleStatusCanceled), sale.Status)
		}
	})
}

func (s *SaleSuite) CreateInvalid(t *testing.T) {
	c := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
	c.Status = string(serviceModels.CarStatusAvailable)
	car, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)

	t.Run("empty car id", func(t *testing.T) {
		sale := models.GenerateSale(car, s.RandomClient().ID)
		sale.CarID = ""
		_, code, err := s.base.client.CreateSale(s.base.ctx, s.adminToken.AccessToken, sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid car id", func(t *testing.T) {
		sale := models.GenerateSale(car, s.RandomClient().ID)
		sale.CarID = "invalid"
		_, code, err := s.base.client.CreateSale(s.base.ctx, s.adminToken.AccessToken, sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing car", func(t *testing.T) {
		sale := models.GenerateSale(car, s.RandomClient().ID)
		sale.CarID = gofakeit.UUID()
		_, code, err := s.base.client.CreateSale(s.base.ctx, s.adminToken.AccessToken, sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty client id", func(t *testing.T) {
		sale := models.GenerateSale(car, "")
		_, code, err := s.base.client.CreateSale(s.base.ctx, s.adminToken.AccessToken, sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid client id", func(t *testing.T) {
		sale := models.GenerateSale(car, "invalid")
		_, code, err := s.base.client.CreateSale(s.base.ctx, s.adminToken.AccessToken, sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing client", func(t *testing.T) {
		sale := models.GenerateSale(car, gofakeit.UUID())
		_, code, err := s.base.client.CreateSale(s.base.ctx, s.adminToken.AccessToken, sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not available car", func(t *testing.T) {
		car := models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID)
		car.Status = string(serviceModels.CarStatusIncoming)
		carInternal, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, car)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)

		sale := models.GenerateSale(carInternal, s.RandomClient().ID)
		_, code, err = s.base.client.CreateSale(s.base.ctx, s.adminToken.AccessToken, sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *SaleSuite) CreateForbidden(t *testing.T) {
	car, code, err := s.base.client.CreateCar(s.base.ctx, s.adminToken.AccessToken, models.GenerateCar(s.RandomModel().ID, s.RandomSupplier().ID))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)

	t.Run("without token", func(t *testing.T) {
		sale := models.GenerateSale(car, s.RandomClient().ID)
		_, code, err := s.base.client.CreateSale(s.base.ctx, "", sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with invalid token", func(t *testing.T) {
		sale := models.GenerateSale(car, s.RandomClient().ID)
		_, code, err := s.base.client.CreateSale(s.base.ctx, "invalid", sale)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})
}

func (s *SaleSuite) Get(t *testing.T) {
	for id := range s.sales {
		sale, code, err := s.base.client.GetSale(s.base.ctx, s.adminToken.AccessToken, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, s.sales[id], sale)
	}
}

func (s *SaleSuite) GetInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, code, err := s.base.client.GetSale(s.base.ctx, s.adminToken.AccessToken, "invalid")
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing id", func(t *testing.T) {
		_, code, err := s.base.client.GetSale(s.base.ctx, s.adminToken.AccessToken, gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *SaleSuite) GetForbidden(t *testing.T) {
	t.Run("without token", func(t *testing.T) {
		_, code, err := s.base.client.GetSale(s.base.ctx, "", s.salesIDs[0])
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with invalid token", func(t *testing.T) {
		_, code, err := s.base.client.GetSale(s.base.ctx, "invalid", s.salesIDs[0])
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})
}

func (s *SaleSuite) UpdateInvalid(t *testing.T) {
	sales := make([]*httpModels.Sale, 1)
	t.Run("create", s.CreateSales(sales))
	sale := sales[0]

	t.Run("invalid id", func(t *testing.T) {
		_, code, err := s.base.client.UpdateSale(s.base.ctx, s.adminToken.AccessToken, "invalid", string(serviceModels.SaleStatusCompleted))
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		_, code, err := s.base.client.UpdateSale(s.base.ctx, s.adminToken.AccessToken, gofakeit.UUID(), string(serviceModels.SaleStatusCompleted))
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid status", func(t *testing.T) {
		_, code, err := s.base.client.UpdateSale(s.base.ctx, s.adminToken.AccessToken, sale.ID, "invalid")
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("update completed sale", func(t *testing.T) {
		sale, code, err := s.base.client.UpdateSale(s.base.ctx, s.adminToken.AccessToken, sale.ID, string(serviceModels.SaleStatusCompleted))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, string(serviceModels.SaleStatusCompleted), sale.Status)
		s.sales[sale.ID] = sale
		_, code, err = s.base.client.UpdateSale(s.base.ctx, s.adminToken.AccessToken, sale.ID, string(serviceModels.SaleStatusCompleted))
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *SaleSuite) UpdateForbidden(t *testing.T) {
	sales := make([]*httpModels.Sale, 1)
	t.Run("create", s.CreateSales(sales))
	sale := sales[0]
	t.Run("without token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateSale(s.base.ctx, "", sale.ID, string(serviceModels.SaleStatusCompleted))
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with invalid token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateSale(s.base.ctx, "invalid", sale.ID, string(serviceModels.SaleStatusCompleted))
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with another manager token", func(t *testing.T) {
		var token string
		for _, manager := range s.employees {
			if manager.ID != sale.EmployeeID {
				token = s.employeesTokens[manager.ID].AccessToken
				break
			}
		}
		require.NotEmpty(t, token, "manager token not found")
		_, code, err := s.base.client.UpdateSale(s.base.ctx, token, sale.ID, string(serviceModels.SaleStatusCompleted))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})
}
