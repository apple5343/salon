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

type SupplierSuite struct {
	suite.Suite
	base         *BaseTestSuite
	adminToken   *models.EmployeeToken
	managerToken *models.EmployeeToken
	suppliers    map[string]*httpModels.Supplier
}

func (s *SupplierSuite) SetupSuite() {
	token, code, err := s.base.client.LoginEmployee(s.base.ctx, s.base.adminCreds.Email, s.base.adminCreds.Password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.adminToken = token
	s.suppliers = make(map[string]*httpModels.Supplier)

	manager := models.GenerateEmployee()
	password := manager.Password
	manager, code, err = s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, manager)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, manager.ID)
	token, code, err = s.base.client.LoginEmployee(s.base.ctx, manager.Email, password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.managerToken = token

	s.T().Run("create", s.Create)
}

func (s *SupplierSuite) TestGet() {
	s.T().Run("get", s.Get)
	s.T().Run("invalid", s.GetInvalid)
}

func (s *SupplierSuite) TestCreate() {
	s.T().Run("create", s.Create)
	s.T().Run("invalid", s.CreateInvalid)
	s.T().Run("forbidden", s.CreateForbidden)
}

func (s *SupplierSuite) TestUpdate() {
	s.T().Run("update", s.Update)
	s.T().Run("get after update", s.Get)
	s.T().Run("invalid", s.UpdateInvalid)
	s.T().Run("forbidden", s.UpdateForbidden)
}

func (s *SupplierSuite) TestList() {
	s.T().Run("list", s.List)
}

func (s *SupplierSuite) CompareSuppliersPublic(t *testing.T, expected *httpModels.Supplier, actual *httpModels.Supplier) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.CountryCode, actual.CountryCode)
	require.Empty(t, actual.CreatedAt)
	require.Empty(t, actual.UpdatedAt)
}

func (s *SupplierSuite) Create(t *testing.T) {
	const suppliersCount int = 10

	for i := 0; i < suppliersCount; i++ {
		supplier, code, err := s.base.client.CreateSupplier(s.base.ctx, s.adminToken.AccessToken, models.GenerateSupplier())
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.suppliers[supplier.ID] = supplier
	}
}

func (s *SupplierSuite) CreateInvalid(t *testing.T) {
	t.Run("invalid country code", func(t *testing.T) {
		supplier := models.GenerateSupplier()
		supplier.CountryCode = gofakeit.Word()
		_, code, err := s.base.client.CreateSupplier(s.base.ctx, s.adminToken.AccessToken, supplier)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty name", func(t *testing.T) {
		supplier := models.GenerateSupplier()
		supplier.Name = ""
		_, code, err := s.base.client.CreateSupplier(s.base.ctx, s.adminToken.AccessToken, supplier)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *SupplierSuite) CreateForbidden(t *testing.T) {
	t.Run("create without token", func(t *testing.T) {
		_, code, err := s.base.client.CreateSupplier(s.base.ctx, "", models.GenerateSupplier())
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("create with manager role", func(t *testing.T) {
		_, code, err := s.base.client.CreateSupplier(s.base.ctx, s.managerToken.AccessToken, models.GenerateSupplier())
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})

}

func (s *SupplierSuite) Get(t *testing.T) {
	t.Run("get internal", func(t *testing.T) {
		for id := range s.suppliers {
			supplier, code, err := s.base.client.GetSupplier(s.base.ctx, s.adminToken.AccessToken, id)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, s.suppliers[id], supplier)
		}
	})

	t.Run("get public", func(t *testing.T) {
		for id := range s.suppliers {
			supplier, code, err := s.base.client.GetSupplier(s.base.ctx, "", id)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			s.CompareSuppliersPublic(t, s.suppliers[id], supplier)
		}
	})
}

func (s *SupplierSuite) GetInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, code, err := s.base.client.GetSupplier(s.base.ctx, s.adminToken.AccessToken, gofakeit.Word())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		_, code, err := s.base.client.GetSupplier(s.base.ctx, s.adminToken.AccessToken, gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *SupplierSuite) Update(t *testing.T) {
	for id, supplier := range s.suppliers {
		newSupplier := models.GenerateSupplier()
		supplier.CountryCode = newSupplier.CountryCode
		supplier.Name = newSupplier.Name
		updated, code, err := s.base.client.UpdateSupplier(s.base.ctx, s.adminToken.AccessToken, supplier)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, supplier.ID, updated.ID)
		require.Equal(t, supplier.Name, updated.Name)
		require.Equal(t, supplier.CountryCode, updated.CountryCode)
		require.Equal(t, supplier.CreatedAt, updated.CreatedAt)
		s.suppliers[id] = updated
	}
}

func (s *SupplierSuite) UpdateInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		supplier := models.GenerateSupplier()
		supplier.ID = gofakeit.Word()
		_, code, err := s.base.client.UpdateSupplier(s.base.ctx, s.adminToken.AccessToken, supplier)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		supplier := models.GenerateSupplier()
		supplier.ID = gofakeit.UUID()
		_, code, err := s.base.client.UpdateSupplier(s.base.ctx, s.adminToken.AccessToken, supplier)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	keys := make([]string, 0, len(s.suppliers))
	for k := range s.suppliers {
		keys = append(keys, k)
	}
	require.GreaterOrEqual(t, len(keys), 1)
	supplier := s.suppliers[keys[0]]

	t.Run("empty name", func(t *testing.T) {
		supplierCopy := models.CopySupplier(supplier)
		supplierCopy.Name = ""
		_, code, err := s.base.client.UpdateSupplier(s.base.ctx, s.adminToken.AccessToken, supplierCopy)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid country code", func(t *testing.T) {
		supplierCopy := models.CopySupplier(supplier)
		supplierCopy.CountryCode = gofakeit.Word()
		_, code, err := s.base.client.UpdateSupplier(s.base.ctx, s.adminToken.AccessToken, supplierCopy)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *SupplierSuite) UpdateForbidden(t *testing.T) {
	t.Run("update without token", func(t *testing.T) {
		for _, supplier := range s.suppliers {
			_, code, err := s.base.client.UpdateSupplier(s.base.ctx, "", supplier)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, code)
			break
		}
	})

	t.Run("update with manager token", func(t *testing.T) {
		for _, supplier := range s.suppliers {
			_, code, err := s.base.client.UpdateSupplier(s.base.ctx, s.managerToken.AccessToken, supplier)
			require.NoError(t, err)
			require.Equal(t, http.StatusForbidden, code)
			break
		}
	})
}

func (s *SupplierSuite) CheckList(t *testing.T, token string, filter *serviceModels.SupplierFilters, expected []string) (map[string]*httpModels.Supplier, []string) {
	all := make(map[string]*httpModels.Supplier)
	sorted := []string{}
	for true {
		suppliers, code, err := s.base.client.GetSuppliers(s.base.ctx, token, filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.LessOrEqual(t, len(suppliers), *filter.Limit)
		for _, sup := range suppliers {
			id := sup.ID
			if _, ok := s.suppliers[id]; !ok {
				t.Errorf("invalid supplier id: %s", id)
				continue
			}
			if _, ok := all[id]; ok {
				t.Errorf("duplicate supplier id: %s", id)
				continue
			}
			all[id] = sup
			sorted = append(sorted, id)
		}
		if len(suppliers) < *filter.Limit {
			break
		}
		*filter.Offset += *filter.Limit
	}

	for _, id := range expected {
		if _, ok := all[id]; !ok {
			t.Errorf("not checked supplier id: %s", id)
		}
	}
	require.Equal(t, len(expected), len(all))

	return all, sorted
}

func (s *SupplierSuite) List(t *testing.T) {
	t.Run("all internal order by created at asc", func(t *testing.T) {
		limit := 2
		offset := 0
		filter := serviceModels.SupplierFilters{
			BaseList: serviceModels.BaseList{
				Limit:  &limit,
				Offset: &offset,
			},
		}
		expected := make([]string, 0, len(s.suppliers))
		for id := range s.suppliers {
			expected = append(expected, id)
		}
		all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
		for i := 1; i < len(sorted); i++ {
			if s.suppliers[sorted[i-1]].CreatedAt.After(s.suppliers[sorted[i]].CreatedAt) {
				t.Errorf("invalid order by created at asc: %v > %v", s.suppliers[sorted[i-1]].CreatedAt, s.suppliers[sorted[i]].CreatedAt)
			}
		}

		for id, sup := range all {
			require.Equal(t, s.suppliers[id], sup)
		}
	})

	t.Run("all public with county code & order by updated at desc", func(t *testing.T) {
		countries := make(map[string][]*httpModels.Supplier)
		for _, s := range s.suppliers {
			countries[s.CountryCode] = append(countries[s.CountryCode], s)
		}
		for country, suppliers := range countries {
			orderBy := serviceModels.SupplierOrderByUpdatedAt
			limit := 3
			offset := 0
			direction := serviceModels.OrderDirectionDESC
			filter := serviceModels.SupplierFilters{
				OrderBy:     &orderBy,
				CountryCode: &country,
				BaseList: serviceModels.BaseList{
					Limit:          &limit,
					Offset:         &offset,
					OrderDirection: &direction,
				},
			}
			expected := make([]string, 0, len(suppliers))
			for _, s := range suppliers {
				expected = append(expected, s.ID)
			}
			all, sorted := s.CheckList(t, "", &filter, expected)
			for i := 1; i < len(all); i++ {
				if s.suppliers[sorted[i-1]].UpdatedAt.Before(s.suppliers[sorted[i]].UpdatedAt) {
					t.Errorf("invalid order by updated at desc: %v < %v", s.suppliers[sorted[i-1]].UpdatedAt, s.suppliers[sorted[i]].UpdatedAt)
				}
			}

			for id, sup := range all {
				s.CompareSuppliersPublic(t, s.suppliers[id], sup)
			}
		}
	})

	t.Run("internal with county code and name & sorted by created at asc", func(t *testing.T) {
		names := []string{"a", "b", "c", "d"}
		countriesWithNames := make(map[string]map[string][]*httpModels.Supplier)
		for _, s := range s.suppliers {
			for _, name := range names {
				if models.MatchesPattern(name, s.Name) {
					if _, ok := countriesWithNames[s.CountryCode]; !ok {
						countriesWithNames[s.CountryCode] = make(map[string][]*httpModels.Supplier)
					}
					countriesWithNames[s.CountryCode][name] = append(countriesWithNames[s.CountryCode][name], s)
				}
			}
		}
		for country, names := range countriesWithNames {
			for name, suppliers := range names {
				orderBy := serviceModels.SupplierOrderByCreatedAt
				limit := 2
				offset := 0
				filter := serviceModels.SupplierFilters{
					OrderBy:     &orderBy,
					CountryCode: &country,
					Name:        &name,
					BaseList: serviceModels.BaseList{
						Limit:  &limit,
						Offset: &offset,
					},
				}
				expected := make([]string, 0, len(suppliers))
				for _, s := range suppliers {
					expected = append(expected, s.ID)
				}
				all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
				for i := 1; i < len(all); i++ {
					if s.suppliers[sorted[i-1]].CreatedAt.After(s.suppliers[sorted[i]].CreatedAt) {
						t.Errorf("invalid order by created at asc: %v > %v", s.suppliers[sorted[i-1]].CreatedAt, s.suppliers[sorted[i]].CreatedAt)
					}
				}

				for id, sup := range all {
					require.Equal(t, s.suppliers[id], sup)
				}
			}
		}
	})
}
