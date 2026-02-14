package integration

import (
	"net/http"
	"salon/tests/integration/models"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	httpModels "salon/internal/transport/http/models"
)

type ClientSuite struct {
	suite.Suite
	base         *BaseTestSuite
	adminToken   *models.EmployeeToken
	managerToken *models.EmployeeToken
	clients      map[string]*httpModels.Client
}

func (s *ClientSuite) SetupSuite() {
	token, code, err := s.base.client.LoginEmployee(s.base.ctx, s.base.adminCreds.Email, s.base.adminCreds.Password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.adminToken = token

	manager := models.GenerateEmployee()
	password := manager.Password
	_, code, err = s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, manager)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	token, code, err = s.base.client.LoginEmployee(s.base.ctx, manager.Email, password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.managerToken = token

	s.clients = make(map[string]*httpModels.Client)
	s.T().Run("register", s.Register)
}

func (s *ClientSuite) TestRegister() {
	s.T().Run("register", s.Register)
	s.T().Run("invalid", s.RegisterInvalid)
	s.T().Run("forbidden", s.RegisterForbidden)
}

func (s *ClientSuite) TestGet() {
	s.T().Run("get", s.Get)
	s.T().Run("update", s.Update)
	s.T().Run("get after update", s.Get)
	s.T().Run("invalid", s.GetInvalid)
	s.T().Run("forbidden", s.GetForbidden)
}

func (s *ClientSuite) TestUpdate() {
	s.T().Run("update", s.Update)
	s.T().Run("invalid", s.UpdateInvalid)
	s.T().Run("forbidden", s.UpdateForbidden)
}

func (s *ClientSuite) Register(t *testing.T) {
	const clientsCount int = 5

	t.Run("register with admin token", func(t *testing.T) {
		for i := 0; i < clientsCount; i++ {
			client, code, err := s.base.client.RegisterClient(s.base.ctx, s.adminToken.AccessToken, models.GenerateClient())
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			s.clients[client.ID] = client
		}
	})

	t.Run("register with manager token", func(t *testing.T) {
		for i := 0; i < clientsCount; i++ {
			client, code, err := s.base.client.RegisterClient(s.base.ctx, s.adminToken.AccessToken, models.GenerateClient())
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			s.clients[client.ID] = client
		}
	})
}

func (s *ClientSuite) RegisterInvalid(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		c := models.GenerateClient()
		c.FullName = ""
		_, code, err := s.base.client.RegisterClient(s.base.ctx, s.managerToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid phone", func(t *testing.T) {
		c := models.GenerateClient()
		c.Phone += "1"
		_, code, err := s.base.client.RegisterClient(s.base.ctx, s.managerToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid email", func(t *testing.T) {
		c := models.GenerateClient()
		c.Email = gofakeit.Word()
		_, code, err := s.base.client.RegisterClient(s.base.ctx, s.managerToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid passport series", func(t *testing.T) {
		c := models.GenerateClient()
		c.Passport.Series += "1"
		_, code, err := s.base.client.RegisterClient(s.base.ctx, s.managerToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid passport number", func(t *testing.T) {
		c := models.GenerateClient()
		c.Passport.Number += "1"
		_, code, err := s.base.client.RegisterClient(s.base.ctx, s.managerToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty passport issued by", func(t *testing.T) {
		c := models.GenerateClient()
		c.Passport.IssuedBy = ""
		_, code, err := s.base.client.RegisterClient(s.base.ctx, s.managerToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("password less that 8", func(t *testing.T) {
		c := models.GenerateClient()
		c.Password = "1234567"
		_, code, err := s.base.client.RegisterClient(s.base.ctx, s.managerToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid birth date", func(t *testing.T) {
		c := models.GenerateClient()
		c.BirthDate = gofakeit.Date().Format("02-01-2006")
		_, code, err := s.base.client.RegisterClient(s.base.ctx, s.managerToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *ClientSuite) RegisterForbidden(t *testing.T) {
	t.Run("register without token", func(t *testing.T) {
		c := models.GenerateClient()
		_, code, err := s.base.client.RegisterClient(s.base.ctx, "", c)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})
}

func (s *ClientSuite) Get(t *testing.T) {
	for id := range s.clients {
		client, code, err := s.base.client.GetClient(s.base.ctx, s.adminToken.AccessToken, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, s.clients[id], client)
	}
}

func (s *ClientSuite) GetInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, code, err := s.base.client.GetClient(s.base.ctx, s.adminToken.AccessToken, gofakeit.Word())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		_, code, err := s.base.client.GetClient(s.base.ctx, s.adminToken.AccessToken, gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *ClientSuite) GetForbidden(t *testing.T) {
	t.Run("get without token", func(t *testing.T) {
		_, code, err := s.base.client.GetClient(s.base.ctx, "", gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})
}

func (s *ClientSuite) Update(t *testing.T) {
	for _, c := range s.clients {
		newc := models.GenerateClient()
		c.FullName = newc.FullName
		c.Phone = newc.Phone
		c.Email = newc.Email
		c.Passport.Series = newc.Passport.Series
		c.Passport.Number = newc.Passport.Number
		c.Passport.IssuedBy = newc.Passport.IssuedBy
		c.BirthDate = newc.BirthDate
		updated, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, c.ID, updated.ID)
		require.Equal(t, c.FullName, updated.FullName)
		require.Equal(t, c.Phone, updated.Phone)
		require.Equal(t, c.Email, updated.Email)
		require.Equal(t, c.Passport.Series, updated.Passport.Series)
		require.Equal(t, c.Passport.Number, updated.Passport.Number)
		require.Equal(t, c.Passport.IssuedBy, updated.Passport.IssuedBy)
		require.Equal(t, c.BirthDate, updated.BirthDate)
		s.clients[c.ID] = updated
	}
}

func (s *ClientSuite) UpdateInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		c := models.GenerateClient()
		c.ID = gofakeit.Word()
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		c := models.GenerateClient()
		c.ID = gofakeit.UUID()
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	require.GreaterOrEqual(t, len(s.clients), 2)
	keys := make([]string, 0, len(s.clients))
	for k := range s.clients {
		keys = append(keys, k)
	}
	c1 := s.clients[keys[0]]
	c2 := s.clients[keys[1]]

	t.Run("empty name", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.FullName = ""
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid phone", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.Phone += "1"
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid email", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.Email = gofakeit.Word()
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid passport series", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.Passport.Series = gofakeit.Word()
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid passport number", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.Passport.Number = gofakeit.Word()
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid passport issued by", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.Passport.IssuedBy = ""
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("existing phone", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.Phone = c2.Phone
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("existing email", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.Email = c2.Email
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("existing passport", func(t *testing.T) {
		c := models.CopyClient(c1)
		c.Passport = c2.Passport
		_, code, err := s.base.client.UpdateClient(s.base.ctx, s.adminToken.AccessToken, c)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})
}

func (s *ClientSuite) UpdateForbidden(t *testing.T) {
	require.GreaterOrEqual(t, len(s.clients), 1)
	var c *httpModels.Client
	for id := range s.clients {
		c = s.clients[id]
		break
	}
	t.Run("update without token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateClient(s.base.ctx, "", c)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("update with invalid token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateClient(s.base.ctx, gofakeit.Word(), c)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})
}
