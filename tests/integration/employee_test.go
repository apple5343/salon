package integration

import (
	"net/http"
	"testing"

	httpModels "salon/internal/transport/http/models"
	"salon/tests/integration/models"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	EmployeeSatausActive   = "active"
	EmployeeSatausInactive = "inactive"
)

type EmployeeSuite struct {
	suite.Suite
	base            *BaseTestSuite
	employees       map[string]*httpModels.Employee
	employeesCreds  map[string]*models.EmployeeCreds
	employeesTokens map[string]*models.EmployeeToken
	adminToken      *models.EmployeeToken
}

func (s *EmployeeSuite) SetupSuite() {
	s.employees = make(map[string]*httpModels.Employee)
	s.employeesCreds = make(map[string]*models.EmployeeCreds)
	s.employeesTokens = make(map[string]*models.EmployeeToken)
	t, code, err := s.base.client.LoginEmployee(s.base.ctx, s.base.adminCreds.Email, s.base.adminCreds.Password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.adminToken = t

	s.T().Run("register", s.Register)
}

func (s *EmployeeSuite) TestGet() {
	s.T().Run("get", s.GetByID)
	s.T().Run("update", s.Update)
	s.T().Run("get after update", s.GetByID)
	s.T().Run("invalid", s.GetByIDInvalid)
	s.T().Run("forbidden", s.GetByIDForbidden)
}

func (s *EmployeeSuite) TestUpdate() {
	s.T().Run("update", s.Update)
	s.T().Run("invalid", s.UpdateInvalid)
	s.T().Run("forbidden", s.UpdateForbidden)
}

func (s *EmployeeSuite) TestAuth() {
	s.T().Run("hire", s.Hire)
	s.T().Run("login", s.Login)
	s.T().Run("update", s.Update)
	s.T().Run("login after update password", s.GetByID)
	s.T().Run("get refresh token", s.GetRefreshToken)
	s.T().Run("get access token", s.GetAccessToken)
	s.T().Run("get profile", s.Profile)
	s.T().Run("register", s.Register)
	s.T().Run("invalid", s.AuthInvalid)
}

func (s *EmployeeSuite) TestRegister() {
	s.T().Run("register", s.Register)
	s.T().Run("invalid", s.RegisterInvalidData)
	s.T().Run("forbidden", s.RegisterForbidden)
}

func (s *EmployeeSuite) TestHire() {
	s.T().Run("hire", s.Hire)
	s.T().Run("invalid", s.HireInvalid)
	s.T().Run("forbidden", s.HireForbidden)
}

func (s *EmployeeSuite) Hire(t *testing.T) {
	for id := range s.employees {
		if s.employees[id].Status == EmployeeSatausActive {
			continue
		}

		e, code, err := s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.employees[id] = e
	}
}

func (s *EmployeeSuite) HireInvalid(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		_, code, err := s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, "invalid")
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("not existing", func(t *testing.T) {
		_, code, err := s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("already hired", func(t *testing.T) {
		found := false
		for id := range s.employees {
			if s.employees[id].Status == EmployeeSatausInactive {
				continue
			}

			_, code, err := s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, id)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, code)
			found = true
			break
		}
		require.True(t, found, "Can't find active employee")
	})
}

func (s *EmployeeSuite) HireForbidden(t *testing.T) {
	var token *models.EmployeeToken
	for id := range s.employees {
		tk, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		token = tk
	}

	require.NotNil(t, token)
	t.Run("manager hire employee", func(t *testing.T) {
		_, code, err := s.base.client.HireEmployee(s.base.ctx, token.AccessToken, gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})

	t.Run("hire without token", func(t *testing.T) {
		_, code, err := s.base.client.HireEmployee(s.base.ctx, "", gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("hire with invalid token", func(t *testing.T) {
		_, code, err := s.base.client.HireEmployee(s.base.ctx, "invalid", gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})
}

func (s *EmployeeSuite) Register(t *testing.T) {
	const employeesCount int = 5
	for i := 0; i < employeesCount; i++ {
		e := models.GenerateEmployee()
		password := e.Password
		e, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Empty(t, e.Password)
		s.employees[e.ID] = e
		s.employeesCreds[e.ID] = &models.EmployeeCreds{
			Email:    e.Email,
			Password: password,
		}
	}
}

func (s *EmployeeSuite) RegisterInvalidData(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.FullName = ""
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty phone", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Phone = ""
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty passport series", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Passport.Series = ""
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty passport number", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Passport.Number = ""
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty passport issued by", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Passport.IssuedBy = ""
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty password", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Password = ""
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty email", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Email = ""
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("empty role", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Role = ""
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *EmployeeSuite) RegisterForbidden(t *testing.T) {
	t.Run("register without token", func(t *testing.T) {
		e := models.GenerateEmployee()
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, "", e)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	var token *models.EmployeeToken
	for id := range s.employees {
		tk, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		token = tk
	}
	require.NotNil(t, token)

	t.Run("manager register employee", func(t *testing.T) {
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, token.AccessToken, models.GenerateEmployee())
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})

	t.Run("admin register admin", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Role = models.RoleAdmin
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *EmployeeSuite) Reregister(t *testing.T) {
	t.Run("existing email", func(t *testing.T) {
		for id := range s.employees {
			e := models.GenerateEmployee()
			e.Email = s.employees[id].Email
			_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
			require.NoError(t, err)
			require.Equal(t, http.StatusConflict, code)
		}
	})

	t.Run("existing phone", func(t *testing.T) {
		for id := range s.employees {
			e := models.GenerateEmployee()
			e.Phone = s.employees[id].Phone
			_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
			require.NoError(t, err)
			require.Equal(t, http.StatusConflict, code)
		}
	})

	t.Run("existing passport", func(t *testing.T) {
		for id := range s.employees {
			e := models.GenerateEmployee()
			e.Passport = s.employees[id].Passport
			_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
			require.NoError(t, err)
			require.Equal(t, http.StatusConflict, code)
		}
	})
}

func (s *EmployeeSuite) GetByID(t *testing.T) {
	for id := range s.employees {
		e, code, err := s.base.client.GetEmployee(s.base.ctx, s.adminToken.AccessToken, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, s.employees[id], e)
	}
}

func (s *EmployeeSuite) GetByIDForbidden(t *testing.T) {
	t.Run("get without token", func(t *testing.T) {
		for id := range s.employees {
			_, code, err := s.base.client.GetEmployee(s.base.ctx, "", id)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, code)
		}
	})

	var token *models.EmployeeToken
	for id := range s.employees {
		tk, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		token = tk
	}
	require.NotNil(t, token)

	t.Run("get with manager role", func(t *testing.T) {
		for id := range s.employees {
			_, code, err := s.base.client.GetEmployee(s.base.ctx, token.AccessToken, id)
			require.NoError(t, err)
			require.Equal(t, http.StatusForbidden, code)
		}
	})
}

func (s *EmployeeSuite) GetByIDInvalid(t *testing.T) {
	t.Run("not existing employee", func(t *testing.T) {
		_, code, err := s.base.client.GetEmployee(s.base.ctx, s.adminToken.AccessToken, gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		_, code, err := s.base.client.GetEmployee(s.base.ctx, s.adminToken.AccessToken, gofakeit.Word())
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *EmployeeSuite) Update(t *testing.T) {
	for id, e := range s.employees {
		newE := models.GenerateEmployee()
		e.FullName = newE.FullName
		e.Phone = newE.Phone
		e.Passport = newE.Passport
		updated, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, e.FullName, updated.FullName)
		require.Equal(t, e.Phone, updated.Phone)
		require.Equal(t, e.Passport, updated.Passport)
		s.employees[id] = updated
	}

	t.Run("update password", func(t *testing.T) {
		for id, e := range s.employees {
			e.Password = gofakeit.Password(true, true, true, true, false, 8)
			password := e.Password
			updated, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			_, code, err = s.base.client.LoginEmployee(s.base.ctx, e.Email, e.Password)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			s.employees[id] = updated
			s.employeesCreds[id].Password = password
		}
	})
}

func (s *EmployeeSuite) UpdateInvalid(t *testing.T) {
	t.Run("not existing employee", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.ID = gofakeit.UUID()
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid id", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.ID = gofakeit.Word()
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	keys := make([]string, 0, len(s.employees))
	for k := range s.employees {
		keys = append(keys, k)
	}
	require.GreaterOrEqual(t, len(keys), 2)
	e1 := s.employees[keys[0]]
	e2 := s.employees[keys[1]]

	t.Run("invalid email", func(t *testing.T) {
		e := models.CopyEmployee(e1)
		e.Email = gofakeit.Word()
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("existing email", func(t *testing.T) {
		e := models.CopyEmployee(e1)
		e.Email = e2.Email
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("invalid phone", func(t *testing.T) {
		e := models.CopyEmployee(e1)
		e.Phone = models.RandomPhone() + "1"
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("existing phone", func(t *testing.T) {
		e := models.CopyEmployee(e1)
		e.Phone = e2.Phone
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("invalid passport series", func(t *testing.T) {
		e := models.CopyEmployee(e1)
		e.Passport = models.RandomPassport()
		e.Passport.Series += "1"
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("invalid passport number", func(t *testing.T) {
		e := models.CopyEmployee(e1)
		e.Passport = models.RandomPassport()
		e.Passport.Number += "1"
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("existing passport", func(t *testing.T) {
		e := models.CopyEmployee(e1)
		e.Passport = e2.Passport
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("empty name", func(t *testing.T) {
		e := models.CopyEmployee(e1)
		e.FullName = ""
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *EmployeeSuite) UpdateForbidden(t *testing.T) {
	t.Run("update without token", func(t *testing.T) {
		for _, e := range s.employees {
			_, code, err := s.base.client.UpdateEmployee(s.base.ctx, "", e)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, code)
		}
	})

	var token *models.EmployeeToken
	for id := range s.employees {
		tk, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		token = tk
		break
	}
	require.NotNil(t, token)

	t.Run("update with manager role", func(t *testing.T) {
		for _, e := range s.employees {
			_, code, err := s.base.client.UpdateEmployee(s.base.ctx, token.AccessToken, e)
			require.NoError(t, err)
			require.Equal(t, http.StatusForbidden, code)
		}
	})
}

func (s *EmployeeSuite) Login(t *testing.T) {
	for id := range s.employees {
		token, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.employeesTokens[id] = token
	}
}

func (s *EmployeeSuite) GetRefreshToken(t *testing.T) {
	for id := range s.employees {
		token, code, err := s.base.client.GetRefreshToken(s.base.ctx, s.employeesTokens[id].RefreshToken)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.employeesTokens[id].RefreshToken = token
	}
}

func (s *EmployeeSuite) GetAccessToken(t *testing.T) {
	for id := range s.employees {
		token, code, err := s.base.client.GetAccessToken(s.base.ctx, s.employeesTokens[id].RefreshToken)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.employeesTokens[id].AccessToken = token
	}
}

func (s *EmployeeSuite) Profile(t *testing.T) {
	for id := range s.employees {
		e, code, err := s.base.client.Profile(s.base.ctx, s.employeesTokens[id].AccessToken)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, s.employees[id], e)
	}
}

func (s *EmployeeSuite) AuthInvalid(t *testing.T) {
	t.Run("login invalid credentials", func(t *testing.T) {
		for id := range s.employees {
			_, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password+"1")
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, code)
		}
	})

	t.Run("login inactive employee", func(t *testing.T) {
		count := 0
		for id := range s.employees {
			if s.employees[id].Status == EmployeeSatausActive {
				continue
			}
			count++
			_, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, code)
		}
		require.Greater(t, count, 0, "Count of inactive employees must be greater than 0")
	})

	t.Run("get refresh token without token", func(t *testing.T) {
		_, code, err := s.base.client.GetRefreshToken(s.base.ctx, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("get refresh toke with invalid refresh token", func(t *testing.T) {
		_, code, err := s.base.client.GetRefreshToken(s.base.ctx, gofakeit.Word())
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("get access token without token", func(t *testing.T) {
		_, code, err := s.base.client.GetAccessToken(s.base.ctx, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("get access token with invalid refresh token", func(t *testing.T) {
		_, code, err := s.base.client.GetAccessToken(s.base.ctx, gofakeit.Word())
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("get profile without token", func(t *testing.T) {
		_, code, err := s.base.client.Profile(s.base.ctx, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("get profile with invalid token", func(t *testing.T) {
		_, code, err := s.base.client.Profile(s.base.ctx, gofakeit.Word())
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})
}
