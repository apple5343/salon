package integration

import (
	"net/http"
	serviceModels "salon/internal/models"
	httpModels "salon/internal/transport/http/models"
	"salon/tests/integration/models"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ModelSuite struct {
	suite.Suite
	base         *BaseTestSuite
	adminToken   *models.EmployeeToken
	managerToken *models.EmployeeToken
	models       map[string]*httpModels.ModelInternalResponse
	modelsIds    []string
	brands       map[string]*httpModels.BrandInternalResponse
	brandsIds    []string
}

func (s *ModelSuite) SetupSuite() {
	s.models = make(map[string]*httpModels.ModelInternalResponse)
	s.brands = make(map[string]*httpModels.BrandInternalResponse)

	s.adminToken = s.base.LoginAdmin(s.T())
	s.managerToken, _ = s.base.HireManager(s.T(), s.adminToken.AccessToken)

	for i := 0; i < 5; i++ {
		brand, code, err := s.base.client.CreateBrand(s.base.ctx, s.adminToken.AccessToken, models.GenerateBrand())
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, code)
		s.brands[brand.ID] = brand
		s.brandsIds = append(s.brandsIds, brand.ID)
	}

	s.T().Run("create", s.Create)
}

func CompareModelsPublic(t *testing.T, expected, actual *httpModels.ModelInternalResponse) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Generation, actual.Generation)
	require.Equal(t, expected.BodyType, actual.BodyType)
	require.Equal(t, expected.TransmissionType, actual.TransmissionType)
	require.Equal(t, expected.FuelType, actual.FuelType)
	require.Equal(t, expected.EngineDisplacement, actual.EngineDisplacement)
	require.Equal(t, expected.PowerHP, actual.PowerHP)
	require.Equal(t, expected.DriveType, actual.DriveType)
	require.Equal(t, expected.BasePrice, actual.BasePrice)
	CompareBrandsPublic(t, expected.Brand, actual.Brand)
	require.Zero(t, actual.CreatedAt)
	require.Zero(t, actual.UpdatedAt)
}

func (s *ModelSuite) RandomModel() *httpModels.Model {
	s.Require().GreaterOrEqual(len(s.modelsIds), 1)
	return models.ModelInternalToModel(s.models[s.modelsIds[gofakeit.Number(0, len(s.modelsIds)-1)]])
}

func (s *ModelSuite) RandomBrandId() string {
	s.Require().GreaterOrEqual(len(s.brandsIds), 1)
	return s.brandsIds[gofakeit.Number(0, len(s.brandsIds)-1)]
}

func (s *ModelSuite) TestCreate() {
	s.T().Run("create", s.Create)
	s.T().Run("invalid", s.CreateInvalid)
	s.T().Run("forbidden", s.CreateForbidden)
}

func (s *ModelSuite) TestGet() {
	s.T().Run("get", s.Get)
	s.T().Run("invalid", s.GetInvalid)
}

func (s *ModelSuite) TestUpdate() {
	s.T().Run("update", s.Update)
	s.T().Run("invalid", s.UpdateInvalid)
	s.T().Run("forbidden", s.UpdateForbidden)
}

func (s *ModelSuite) Create(t *testing.T) {
	const modelsCount int = 10

	for i := 0; i < modelsCount; i++ {
		m := models.GenerateModel(s.RandomBrandId())
		model, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, m)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.models[model.ID] = model
		s.modelsIds = append(s.modelsIds, model.ID)
	}
}

func (s *ModelSuite) CreateInvalid(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.Name = ""
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid brand id", func(t *testing.T) {
		model := models.GenerateModel("invalid")
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing brand", func(t *testing.T) {
		model := models.GenerateModel(gofakeit.UUID())
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty brand id", func(t *testing.T) {
		model := models.GenerateModel("")
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty generation", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.Generation = ""
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty body type", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.BodyType = ""
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty transmission type", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.TransmissionType = ""
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty fuel type", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.Name = ""
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("negative engine displacement", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.EngineDisplacement = -1
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("negative power hp", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.PowerHP = -1
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty drive type", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.DriveType = ""
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("negative base price", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.BasePrice = -1
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid body type", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.BodyType = "invalid"
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid transmission type", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.TransmissionType = "invalid"
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid fuel type", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.FuelType = "invalid"
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid drive type", func(t *testing.T) {
		model := models.GenerateModel(s.RandomBrandId())
		model.DriveType = "invalid"
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.adminToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *ModelSuite) CreateForbidden(t *testing.T) {
	t.Run("without token", func(t *testing.T) {
		_, code, err := s.base.client.CreateModel(s.base.ctx, "", models.GenerateModel(s.RandomBrandId()))
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, code, err := s.base.client.CreateModel(s.base.ctx, gofakeit.Word(), models.GenerateModel(s.RandomBrandId()))
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with manager token", func(t *testing.T) {
		_, code, err := s.base.client.CreateModel(s.base.ctx, s.managerToken.AccessToken, models.GenerateModel(s.RandomBrandId()))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})
}

func (s *ModelSuite) Get(t *testing.T) {
	t.Run("get internal", func(t *testing.T) {
		for id := range s.models {
			model, code, err := s.base.client.GetModel(s.base.ctx, s.adminToken.AccessToken, id)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, s.models[id], model)
		}
	})

	t.Run("get public", func(t *testing.T) {
		for id := range s.models {
			model, code, err := s.base.client.GetModel(s.base.ctx, "", id)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			CompareModelsPublic(t, s.models[id], model)
		}
	})
}

func (s *ModelSuite) GetInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, code, err := s.base.client.GetModel(s.base.ctx, s.adminToken.AccessToken, gofakeit.Word())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		_, code, err := s.base.client.GetModel(s.base.ctx, s.adminToken.AccessToken, gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *ModelSuite) Update(t *testing.T) {
	for id, model := range s.models {
		newModel := models.GenerateModel(s.RandomBrandId())
		newModel.ID = id
		newBrand := s.brands[newModel.BrandID]
		updated, code, err := s.base.client.UpdateModel(s.base.ctx, s.adminToken.AccessToken, newModel)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, newModel.ID, updated.ID)
		require.Equal(t, newModel.Name, updated.Name)
		require.Equal(t, newModel.Generation, updated.Generation)
		require.Equal(t, newModel.BodyType, updated.BodyType)
		require.Equal(t, newModel.TransmissionType, updated.TransmissionType)
		require.Equal(t, newModel.FuelType, updated.FuelType)
		require.Equal(t, newModel.EngineDisplacement, updated.EngineDisplacement)
		require.Equal(t, newModel.PowerHP, updated.PowerHP)
		require.Equal(t, newModel.DriveType, updated.DriveType)
		require.Equal(t, newModel.BasePrice, updated.BasePrice)
		require.Equal(t, newBrand.ID, updated.Brand.ID)
		require.Equal(t, newBrand.Name, updated.Brand.Name)
		require.Equal(t, newBrand.Description, updated.Brand.Description)
		require.Equal(t, newBrand.CountryCode, updated.Brand.CountryCode)
		require.Equal(t, model.CreatedAt, updated.CreatedAt)
		require.NotEqual(t, model.UpdatedAt, updated.UpdatedAt)
		s.models[id] = updated
	}
}

func (s *ModelSuite) UpdateInvalid(t *testing.T) {
	type getModelFunc func() *httpModels.Model

	tests := []struct {
		name string
		get  getModelFunc
	}{
		{
			name: "invalid id",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.ID = gofakeit.Word()
				return model
			},
		},
		{
			name: "not existing",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.ID = gofakeit.UUID()
				return model
			},
		},
		{
			name: "invalid brand id",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.BrandID = gofakeit.Word()
				return model
			},
		},
		{
			name: "not existing brand",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.BrandID = gofakeit.UUID()
				return model
			},
		},
		{
			name: "empty name",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.Name = ""
				return model
			},
		},
		{
			name: "empty generation",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.Generation = ""
				return model
			},
		},
		{
			name: "invalid body type",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.BodyType = "invalid"
				return model
			},
		},
		{
			name: "invalid transmission type",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.TransmissionType = "invalid"
				return model
			},
		},
		{
			name: "invalid fuel type",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.FuelType = "invalid"
				return model
			},
		},
		{
			name: "invalid drive type",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.DriveType = "invalid"
				return model
			},
		},
		{
			name: "negative engine displacement",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.EngineDisplacement = -1
				return model
			},
		},
		{
			name: "negative power hp",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.PowerHP = -1
				return model
			},
		},
		{
			name: "negative base price",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.BasePrice = -1
				return model
			},
		},
		{
			name: "empty brand id",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.BrandID = ""
				return model
			},
		},
		{
			name: "empty transmission type",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.TransmissionType = ""
				return model
			},
		},
		{
			name: "empty fuel type",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.FuelType = ""
				return model
			},
		},
		{
			name: "empty drive type",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.DriveType = ""
				return model
			},
		},
		{
			name: "empty body type",
			get: func() *httpModels.Model {
				model := s.RandomModel()
				model.BodyType = ""
				return model
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.get()
			_, code, err := s.base.client.UpdateModel(s.base.ctx, s.adminToken.AccessToken, model)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, code)
		})
	}
}

func (s *ModelSuite) UpdateForbidden(t *testing.T) {
	t.Run("without token", func(t *testing.T) {
		model := s.RandomModel()
		_, code, err := s.base.client.UpdateModel(s.base.ctx, "", model)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with invalid token", func(t *testing.T) {
		model := s.RandomModel()
		_, code, err := s.base.client.UpdateModel(s.base.ctx, gofakeit.Word(), model)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("with manager token", func(t *testing.T) {
		model := s.RandomModel()
		_, code, err := s.base.client.UpdateModel(s.base.ctx, s.managerToken.AccessToken, model)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})
}

func (s *ModelSuite) CheckList(t *testing.T, token string, filter *serviceModels.ModelFilters, expected []string) (map[string]*httpModels.ModelShort, []string) {
	all := make(map[string]*httpModels.ModelShort)
	sorted := []string{}
	for true {
		models, code, err := s.base.client.GetModels(s.base.ctx, token, filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.LessOrEqual(t, len(models), *filter.Limit)
		for _, model := range models {
			id := model.ID
			if _, ok := all[id]; !ok {
				t.Errorf("invalid model id: %s", id)
				continue
			}
			if _, ok := all[id]; ok {
				t.Errorf("duplicate model id: %s", id)
				continue
			}
			all[id] = model
			sorted = append(sorted, id)
		}
		if len(models) < *filter.Limit {
			break
		}
		*filter.Offset += *filter.Limit
	}

	for _, id := range expected {
		if _, ok := all[id]; !ok {
			t.Errorf("not checked model id: %s", id)
		}
	}
	require.Equal(t, len(expected), len(all))
	return all, sorted
}
