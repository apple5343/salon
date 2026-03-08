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

type BrandSuite struct {
	suite.Suite
	base         *BaseTestSuite
	adminToken   *models.EmployeeToken
	managerToken *models.EmployeeToken
	brands       map[string]*httpModels.Brand
}

func (s *BrandSuite) SetupSuite() {
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

	s.brands = make(map[string]*httpModels.Brand)

	s.T().Run("create", s.Create(5))
}

func (s *BrandSuite) TestCreate() {
	s.T().Run("create", s.Create(5))
	s.T().Run("invalid", s.CreateInvalid)
	s.T().Run("forbidden", s.CreateForbidden)
}

func (s *BrandSuite) TestGet() {
	s.T().Run("get", s.Get)
	s.T().Run("update", s.Update)
	s.T().Run("get after update", s.Get)
	s.T().Run("invalid", s.GetInvalid)
}

func (s *BrandSuite) TestUpdate() {
	s.T().Run("update", s.Update)
	s.T().Run("invalid", s.UpdateInvalid)
	s.T().Run("forbidden", s.UpdateForbidden)
}

func (s *BrandSuite) TestList() {
	s.T().Run("create", s.Create(15))
	s.T().Run("list", s.List)
}

func (s *BrandSuite) CompareBrandsPublic(t *testing.T, expected, actual *httpModels.Brand) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Description, actual.Description)
	require.Empty(t, actual.CreatedAt)
	require.Empty(t, actual.UpdatedAt)
}

func (s *BrandSuite) Create(count int) func(t *testing.T) {
	return func(t *testing.T) {
		for i := 0; i < count; i++ {
			brand, code, err := s.base.client.CreateBrand(s.base.ctx, s.adminToken.AccessToken, models.GenerateBrand())
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.brands[brand.ID] = brand
		}
	}
}

func (s *BrandSuite) CreateInvalid(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		brand := models.GenerateBrand()
		brand.Name = ""
		_, code, err := s.base.client.CreateBrand(s.base.ctx, s.adminToken.AccessToken, brand)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	t.Run("empty description", func(t *testing.T) {
		brand := models.GenerateBrand()
		brand.Description = ""
		_, code, err := s.base.client.CreateBrand(s.base.ctx, s.adminToken.AccessToken, brand)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	t.Run("invalid country code", func(t *testing.T) {
		brand := models.GenerateBrand()
		brand.CountryCode = gofakeit.Word()
		_, code, err := s.base.client.CreateBrand(s.base.ctx, s.adminToken.AccessToken, brand)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})
}

func (s *BrandSuite) CreateForbidden(t *testing.T) {
	t.Run("create with manager role", func(t *testing.T) {
		_, code, err := s.base.client.CreateBrand(s.base.ctx, s.managerToken.AccessToken, models.GenerateBrand())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusForbidden, code)
	})

	t.Run("create without token", func(t *testing.T) {
		_, code, err := s.base.client.CreateBrand(s.base.ctx, "", models.GenerateBrand())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusUnauthorized, code)
	})
}

func (s *BrandSuite) Get(t *testing.T) {
	t.Run("get internal", func(t *testing.T) {
		for id := range s.brands {
			brand, code, err := s.base.client.GetBrand(s.base.ctx, s.adminToken.AccessToken, id)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.Require().Equal(s.brands[id], brand)
		}
	})

	t.Run("get public", func(t *testing.T) {
		for id := range s.brands {
			brand, code, err := s.base.client.GetBrand(s.base.ctx, "", id)
			s.Require().NoError(err)
			s.Require().Equal(http.StatusOK, code)
			s.CompareBrandsPublic(t, s.brands[id], brand)
		}
	})
}

func (s *BrandSuite) GetInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, code, err := s.base.client.GetBrand(s.base.ctx, s.adminToken.AccessToken, gofakeit.Word())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		_, code, err := s.base.client.GetBrand(s.base.ctx, s.adminToken.AccessToken, gofakeit.UUID())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})
}

func (s *BrandSuite) Update(t *testing.T) {
	for id, brand := range s.brands {
		newBrand := models.GenerateBrand()
		brand.CountryCode = newBrand.CountryCode
		brand.Description = newBrand.Description
		brand.Name = newBrand.Name
		updated, code, err := s.base.client.UpdateBrand(s.base.ctx, s.adminToken.AccessToken, brand)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		require.Equal(t, brand.Name, updated.Name)
		require.Equal(t, brand.Description, updated.Description)
		require.Equal(t, brand.CountryCode, updated.CountryCode)
		require.Equal(t, brand.CreatedAt, updated.CreatedAt)
		require.NotEqual(t, brand.UpdatedAt, updated.UpdatedAt)
		s.brands[id] = updated
	}
}

func (s *BrandSuite) UpdateInvalid(t *testing.T) {
	var brand *httpModels.Brand
	for id := range s.brands {
		brand = s.brands[id]
		break
	}
	t.Run("empty name", func(t *testing.T) {
		b := models.CopyBrand(brand)
		b.Name = ""
		_, code, err := s.base.client.UpdateBrand(s.base.ctx, s.adminToken.AccessToken, b)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	t.Run("invalid country code", func(t *testing.T) {
		b := models.CopyBrand(brand)
		b.CountryCode = gofakeit.Word()
		_, code, err := s.base.client.UpdateBrand(s.base.ctx, s.adminToken.AccessToken, b)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	t.Run("empty description", func(t *testing.T) {
		b := models.CopyBrand(brand)
		b.Description = ""
		_, code, err := s.base.client.UpdateBrand(s.base.ctx, s.adminToken.AccessToken, b)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		b := models.GenerateBrand()
		b.ID = gofakeit.UUID()
		_, code, err := s.base.client.UpdateBrand(s.base.ctx, s.adminToken.AccessToken, b)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})

	t.Run("invalid id", func(t *testing.T) {
		b := models.GenerateBrand()
		b.ID = gofakeit.Word()
		_, code, err := s.base.client.UpdateBrand(s.base.ctx, s.adminToken.AccessToken, b)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, code)
	})
}

func (s *BrandSuite) UpdateForbidden(t *testing.T) {
	var brand *httpModels.Brand
	for id := range s.brands {
		brand = s.brands[id]
		break
	}
	t.Run("update with manager role", func(t *testing.T) {
		_, code, err := s.base.client.UpdateBrand(s.base.ctx, s.managerToken.AccessToken, brand)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusForbidden, code)
	})

	t.Run("update without token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateBrand(s.base.ctx, "", brand)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusUnauthorized, code)
	})
}

func (s *BrandSuite) CheckList(t *testing.T, token string, filter *serviceModels.BrandFilters, expected []string) (map[string]*httpModels.Brand, []string) {
	all := make(map[string]*httpModels.Brand)
	sorted := []string{}
	for true {
		brands, code, err := s.base.client.GetBrands(s.base.ctx, token, filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.LessOrEqual(t, len(brands), *filter.Limit)
		for _, b := range brands {
			id := b.ID
			if _, ok := s.brands[id]; !ok {
				t.Errorf("invalid brand id: %s", id)
				continue
			}
			if _, ok := all[id]; ok {
				t.Errorf("duplicate brand id: %s", id)
				continue
			}
			all[id] = b
			sorted = append(sorted, id)
		}
		if len(brands) < *filter.Limit {
			break
		}
		*filter.Offset += *filter.Limit
	}

	for _, id := range expected {
		if _, ok := all[id]; !ok {
			t.Errorf("not checked brand id: %s", id)
		}
	}
	require.Equal(t, len(expected), len(all))
	return all, sorted
}

func (s *BrandSuite) List(t *testing.T) {
	t.Run("all internal order by created at asc", func(t *testing.T) {
		// сортировка по created_at asc должна применяться по умолчанию
		limit := 4
		offset := 0
		filter := serviceModels.BrandFilters{
			BaseList: serviceModels.BaseList{
				Limit:  &limit,
				Offset: &offset,
			},
		}
		expected := make([]string, 0, len(s.brands))
		for id := range s.brands {
			expected = append(expected, id)
		}
		all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
		for i := 1; i < len(sorted); i++ {
			if s.brands[sorted[i-1]].CreatedAt.After(s.brands[sorted[i]].CreatedAt) {
				t.Errorf("invalid order by created at asc: %v > %v", s.brands[sorted[i-1]].CreatedAt, s.brands[sorted[i]].CreatedAt)
			}
		}
		for id, b := range all {
			require.Equal(t, s.brands[id], b)
		}
	})

	t.Run("all public order by created at desc", func(t *testing.T) {
		// сортировка по created_at должна применяться по умолчанию
		limit := 4
		offset := 0
		direction := serviceModels.OrderDirectionDESC
		filter := serviceModels.BrandFilters{
			BaseList: serviceModels.BaseList{
				Limit:          &limit,
				Offset:         &offset,
				OrderDirection: &direction,
			},
		}
		expected := make([]string, 0, len(s.brands))
		for id := range s.brands {
			expected = append(expected, id)
		}
		all, sorted := s.CheckList(t, "", &filter, expected)
		for i := 1; i < len(sorted); i++ {
			if s.brands[sorted[i-1]].CreatedAt.Before(s.brands[sorted[i]].CreatedAt) {
				t.Errorf("invalid order by created at desc: %v < %v", s.brands[sorted[i-1]].CreatedAt, s.brands[sorted[i]].CreatedAt)
			}
		}
		for id, b := range all {
			s.CompareBrandsPublic(t, s.brands[id], b)
		}
	})

	t.Run("all internal with country code & order by updated at desc", func(t *testing.T) {
		countries := make(map[string][]*httpModels.Brand)
		for _, b := range s.brands {
			countries[b.CountryCode] = append(countries[b.CountryCode], b)
		}
		for countryCode, bs := range countries {
			orderBy := serviceModels.BrandOrderByUpdatedAt
			limit := 4
			offset := 0
			direction := serviceModels.OrderDirectionDESC
			filter := serviceModels.BrandFilters{
				CountryCode: &countryCode,
				OrderBy:     &orderBy,
				BaseList: serviceModels.BaseList{
					Limit:          &limit,
					Offset:         &offset,
					OrderDirection: &direction,
				},
			}
			expected := make([]string, 0, len(bs))
			for _, b := range bs {
				expected = append(expected, b.ID)
			}
			all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
			for i := 1; i < len(sorted); i++ {
				if s.brands[sorted[i-1]].UpdatedAt.Before(s.brands[sorted[i]].UpdatedAt) {
					t.Errorf("invalid order by updated at desc: %v < %v", s.brands[sorted[i-1]].UpdatedAt, s.brands[sorted[i]].UpdatedAt)
				}
			}
			for id, b := range all {
				require.Equal(t, s.brands[id], b)
			}
		}
	})

	t.Run("public with country code and name & order by created at asc", func(t *testing.T) {
		names := []string{"a", "b", "c", "d"}
		countriesWithNames := make(map[string]map[string][]*httpModels.Brand)
		for _, b := range s.brands {
			for _, name := range names {
				if models.MatchesPattern(name, b.Name) {
					if _, ok := countriesWithNames[b.CountryCode]; !ok {
						countriesWithNames[b.CountryCode] = make(map[string][]*httpModels.Brand)
					}
					countriesWithNames[b.CountryCode][name] = append(countriesWithNames[b.CountryCode][name], b)
				}
			}
		}
		for countryCode, namesWithBrands := range countriesWithNames {
			// сортировка по created_at asc должна применяться по умолчанию
			for name, bs := range namesWithBrands {
				limit := 4
				offset := 0
				filter := serviceModels.BrandFilters{
					CountryCode: &countryCode,
					Name:        &name,
					BaseList: serviceModels.BaseList{
						Limit:  &limit,
						Offset: &offset,
					},
				}
				expected := make([]string, 0, len(bs))
				for _, b := range bs {
					expected = append(expected, b.ID)
				}
				all, sorted := s.CheckList(t, "", &filter, expected)
				for i := 1; i < len(sorted); i++ {
					if s.brands[sorted[i-1]].CreatedAt.After(s.brands[sorted[i]].CreatedAt) {
						t.Errorf("invalid order by created at asc: %v > %v", s.brands[sorted[i-1]].CreatedAt, s.brands[sorted[i]].CreatedAt)
					}
				}
				for id, b := range all {
					s.CompareBrandsPublic(t, s.brands[id], b)
				}
			}
		}
	})
}
