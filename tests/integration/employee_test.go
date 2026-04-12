package integration

import (
	"net/http"
	"strings"
	"testing"

	serviceModels "salon/internal/models"
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
	base              *BaseTestSuite
	employees         map[string]*httpModels.Employee
	employeesIncative map[string]bool
	employeesActive   map[string]bool
	employeesCreds    map[string]*models.EmployeeCreds
	employeesTokens   map[string]*models.EmployeeToken
	adminToken        *models.EmployeeToken
}

func (s *EmployeeSuite) SetupSuite() {
	s.employees = make(map[string]*httpModels.Employee)
	s.employeesIncative = make(map[string]bool)
	s.employeesActive = make(map[string]bool)
	s.employeesCreds = make(map[string]*models.EmployeeCreds)
	s.employeesTokens = make(map[string]*models.EmployeeToken)
	t, code, err := s.base.client.LoginEmployee(s.base.ctx, s.base.adminCreds.Email, s.base.adminCreds.Password)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.adminToken = t

	admin, code, err := s.base.client.Profile(s.base.ctx, s.adminToken.AccessToken)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, code)
	s.employees[admin.ID] = admin
	s.employeesActive[admin.ID] = true
	s.employeesCreds[admin.ID] = &models.EmployeeCreds{
		Email:    s.base.adminCreds.Email,
		Password: s.base.adminCreds.Password,
	}
	s.employeesTokens[admin.ID] = t

	s.T().Run("register", s.Register(10))
	s.T().Run("hire", s.Hire)
}

func (s *EmployeeSuite) LoginEmployee(t *testing.T) *models.EmployeeToken {
	var token *models.EmployeeToken
	for id := range s.employeesActive {
		if s.employees[id].Role == string(serviceModels.EmployeeRoleAdmin) {
			continue
		}
		tk, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		token = tk
		break
	}
	require.NotNil(t, token, "cannot find active employee")
	return token
}

func (s *EmployeeSuite) RandomEmployee(t *testing.T) *httpModels.Employee {
	var e *httpModels.Employee
	for id := range s.employees {
		e = s.employees[id]
		break
	}
	require.NotNil(t, e, "cannot find active employee")
	return e
}

func (s *EmployeeSuite) RandomInactiveEmployee(t *testing.T) *httpModels.Employee {
	var e *httpModels.Employee
	for id := range s.employeesIncative {
		e = s.employees[id]
		break
	}
	require.NotNil(t, e, "cannot find inactive employee")
	return e
}

func (s *EmployeeSuite) RandomActiveEmployee(t *testing.T) *httpModels.Employee {
	var e *httpModels.Employee
	for id := range s.employeesActive {
		e = s.employees[id]
		break
	}
	require.NotNil(t, e, "cannot find active employee")
	return e
}

func (s *EmployeeSuite) TestGet() {
	s.T().Run("get", s.Get)
	s.T().Run("update", s.Update)
	s.T().Run("get after update", s.Get)
	s.T().Run("invalid", s.GetInvalid)
	s.T().Run("forbidden", s.GetForbidden)
}

func (s *EmployeeSuite) TestUpdate() {
	s.T().Run("update", s.Update)
	s.T().Run("invalid", s.UpdateInvalid)
	s.T().Run("forbidden", s.UpdateForbidden)
}

func (s *EmployeeSuite) TestAuth() {
	s.T().Run("register", s.Register(10))
	s.T().Run("login", s.Login)
	s.T().Run("get refresh token", s.GetRefreshToken)
	s.T().Run("get access token", s.GetAccessToken)
	s.T().Run("get profile", s.Profile)
	s.T().Run("register", s.Register(10))
	s.T().Run("invalid", s.AuthInvalid)
}

func (s *EmployeeSuite) TestList() {
	s.T().Run("register", s.Register(15))
	s.T().Run("hire", s.Hire)
	s.T().Run("list", s.List)
	s.T().Run("invalid", s.ListInvalid)
	s.T().Run("forbidden", s.ListForbidden)
}

func (s *EmployeeSuite) TestRegister() {
	s.T().Run("register", s.Register(10))
	s.T().Run("invalid", s.RegisterInvalid)
	s.T().Run("forbidden", s.RegisterForbidden)
}

func (s *EmployeeSuite) Hire(t *testing.T) {
	for id := range s.employeesIncative {
		e, code, err := s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, EmployeeSatausActive, e.Status)
		delete(s.employeesIncative, e.ID)
		s.employeesActive[e.ID] = true
		s.employees[e.ID] = e
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
		hiredEmployee := s.RandomActiveEmployee(t)
		_, code, err := s.base.client.HireEmployee(s.base.ctx, s.adminToken.AccessToken, hiredEmployee.ID)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *EmployeeSuite) HireForbidden(t *testing.T) {
	token := s.LoginEmployee(t)
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

func (s *EmployeeSuite) Register(count int) func(t *testing.T) {
	return func(t *testing.T) {
		for i := 0; i < count; i++ {
			e := models.GenerateEmployee()
			password := e.Password
			e, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			require.Empty(t, e.Password)
			s.employees[e.ID] = e
			s.employeesIncative[e.ID] = true
			s.employeesCreds[e.ID] = &models.EmployeeCreds{
				Email:    e.Email,
				Password: password,
			}
		}
	}
}

func (s *EmployeeSuite) RegisterInvalid(t *testing.T) {
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

	existingEmployee := s.RandomEmployee(t)

	t.Run("existing email", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Email = existingEmployee.Email
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("existing phone", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Phone = existingEmployee.Phone
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})

	t.Run("existing passport", func(t *testing.T) {
		e := models.GenerateEmployee()
		e.Passport = existingEmployee.Passport
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, s.adminToken.AccessToken, e)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, code)
	})
}

func (s *EmployeeSuite) RegisterForbidden(t *testing.T) {
	t.Run("register without token", func(t *testing.T) {
		e := models.GenerateEmployee()
		_, code, err := s.base.client.RegisterEmployee(s.base.ctx, "", e)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	token := s.LoginEmployee(t)

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

func (s *EmployeeSuite) Get(t *testing.T) {
	for id := range s.employees {
		e, code, err := s.base.client.GetEmployee(s.base.ctx, s.adminToken.AccessToken, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, s.employees[id], e)
	}
}

func (s *EmployeeSuite) GetForbidden(t *testing.T) {
	id := s.RandomEmployee(t).ID
	t.Run("get without token", func(t *testing.T) {
		_, code, err := s.base.client.GetEmployee(s.base.ctx, "", id)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	token := s.LoginEmployee(t)

	t.Run("get with manager role", func(t *testing.T) {
		_, code, err := s.base.client.GetEmployee(s.base.ctx, token.AccessToken, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})
}

func (s *EmployeeSuite) GetInvalid(t *testing.T) {
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
	for id := range s.employees {
		if s.employees[id].Role == string(serviceModels.EmployeeRoleAdmin) {
			continue
		}
		e := models.CopyEmployee(s.employees[id])
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
		require.Equal(t, s.employees[id].CreatedAt, updated.CreatedAt)
		require.NotEqual(t, s.employees[id].UpdatedAt, updated.UpdatedAt)
		s.employees[id] = updated
	}

	t.Run("update password", func(t *testing.T) {
		for id := range s.employees {
			if s.employees[id].Role == string(serviceModels.EmployeeRoleAdmin) {
				continue
			}
			e := models.CopyEmployee(s.employees[id])
			e.Password = gofakeit.Password(true, true, true, true, false, 8)
			password := e.Password
			updated, code, err := s.base.client.UpdateEmployee(s.base.ctx, s.adminToken.AccessToken, e)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, code)
			require.Equal(t, s.employees[id].CreatedAt, updated.CreatedAt)
			require.NotEqual(t, s.employees[id].UpdatedAt, updated.UpdatedAt)
			if e.Status == string(serviceModels.EmployeeStatusActive) {
				_, code, err = s.base.client.LoginEmployee(s.base.ctx, e.Email, e.Password)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, code)
			}
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
	existingEmployee := s.RandomEmployee(t)
	t.Run("update without token", func(t *testing.T) {
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, "", existingEmployee)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	token := s.LoginEmployee(t)

	t.Run("update with manager role", func(t *testing.T) {
		_, code, err := s.base.client.UpdateEmployee(s.base.ctx, token.AccessToken, existingEmployee)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})
}

func (s *EmployeeSuite) Login(t *testing.T) {
	for id := range s.employeesActive {
		token, code, err := s.base.client.LoginEmployee(s.base.ctx, s.employeesCreds[id].Email, s.employeesCreds[id].Password)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.employeesTokens[id] = token
	}
}

func (s *EmployeeSuite) GetRefreshToken(t *testing.T) {
	for id := range s.employeesActive {
		token, code, err := s.base.client.GetRefreshToken(s.base.ctx, s.employeesTokens[id].RefreshToken)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.employeesTokens[id].RefreshToken = token
	}
}

func (s *EmployeeSuite) GetAccessToken(t *testing.T) {
	for id := range s.employeesActive {
		token, code, err := s.base.client.GetAccessToken(s.base.ctx, s.employeesTokens[id].RefreshToken)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		s.employeesTokens[id].AccessToken = token
	}
}

func (s *EmployeeSuite) Profile(t *testing.T) {
	for id := range s.employeesActive {
		e, code, err := s.base.client.Profile(s.base.ctx, s.employeesTokens[id].AccessToken)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, s.employees[id], e)
	}
}

func (s *EmployeeSuite) AuthInvalid(t *testing.T) {
	inactiveEmployee := s.RandomInactiveEmployee(t)
	creds := s.employeesCreds[inactiveEmployee.ID]
	t.Run("login invalid credentials", func(t *testing.T) {
		_, code, err := s.base.client.LoginEmployee(s.base.ctx, creds.Email, creds.Password+"1")
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("login inactive employee", func(t *testing.T) {
		_, code, err := s.base.client.LoginEmployee(s.base.ctx, creds.Email, creds.Password)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
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

func (s *EmployeeSuite) CheckList(t *testing.T, token string, filter *serviceModels.EmployeeFilters, expected []string) (map[string]*httpModels.Employee, []string) {
	all := make(map[string]*httpModels.Employee)
	sorted := []string{}
	for true {
		employees, code, err := s.base.client.GetEmployees(s.base.ctx, token, filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, code)
		require.LessOrEqual(t, len(employees), *filter.Limit)
		for _, e := range employees {
			id := e.ID
			if _, ok := all[id]; ok {
				t.Errorf("duplicate employee id: %s", id)
				continue
			}
			all[id] = e
			sorted = append(sorted, id)
		}
		if len(employees) < *filter.Limit {
			break
		}
		*filter.Offset += *filter.Limit
	}

	for _, id := range expected {
		if _, ok := all[id]; !ok {
			t.Errorf("not checked employee id: %s", id)
		}
	}
	require.Equal(t, len(expected), len(all))
	return all, sorted
}

func (s *EmployeeSuite) List(t *testing.T) {
	t.Run("all employees order by created at asc", func(t *testing.T) {
		limit := 4
		offset := 0
		filter := serviceModels.EmployeeFilters{
			BaseList: serviceModels.BaseList{
				Limit:  &limit,
				Offset: &offset,
			},
		}
		expected := make([]string, 0, len(s.employees))
		for id := range s.employees {
			expected = append(expected, id)
		}
		all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
		for i := 1; i < len(sorted); i++ {
			if s.employees[sorted[i-1]].CreatedAt.After(s.employees[sorted[i]].CreatedAt) {
				t.Errorf("invalid order by created at asc: %v > %v", s.employees[sorted[i-1]].CreatedAt, s.employees[sorted[i]].CreatedAt)
			}
		}
		for id, e := range all {
			require.Equal(t, s.employees[id], e)
		}
	})

	t.Run("all employees order by hire date desc", func(t *testing.T) {
		limit := 4
		offset := 0
		direction := serviceModels.OrderDirectionDESC
		orderBy := serviceModels.EmployeeOrderByHireDate
		filter := serviceModels.EmployeeFilters{
			BaseList: serviceModels.BaseList{
				Limit:          &limit,
				Offset:         &offset,
				OrderDirection: &direction,
			},
			OrderBy: &orderBy,
		}
		expected := make([]string, 0, len(s.employees))
		for id := range s.employees {
			expected = append(expected, id)
		}
		all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)

		// Проверяем сортировку по hire_date DESC (NULLS LAST)
		for i := 1; i < len(sorted); i++ {
			prev := s.employees[sorted[i-1]]
			curr := s.employees[sorted[i]]

			// Если у предыдущего есть hire_date, а у текущего нет - это правильно (NULLS LAST)
			if prev.HireDate != nil && curr.HireDate == nil {
				continue
			}

			// Если у обоих есть hire_date, проверяем порядок DESC
			if prev.HireDate != nil && curr.HireDate != nil {
				if prev.HireDate.Before(*curr.HireDate) {
					t.Errorf("invalid order by hire date desc: %v < %v", prev.HireDate, curr.HireDate)
				}
			}

			// Если у обоих nil, проверяем вторичную сортировку по created_at DESC
			if prev.HireDate == nil && curr.HireDate == nil {
				if prev.CreatedAt.Before(curr.CreatedAt) {
					t.Errorf("invalid secondary order by created_at desc: %v < %v", prev.CreatedAt, curr.CreatedAt)
				}
			}
		}

		for id, e := range all {
			require.Equal(t, s.employees[id], e)
		}
	})

	t.Run("active employees only", func(t *testing.T) {
		limit := 4
		offset := 0
		status := serviceModels.EmployeeStatusActive
		filter := serviceModels.EmployeeFilters{
			Status: &status,
			BaseList: serviceModels.BaseList{
				Limit:  &limit,
				Offset: &offset,
			},
		}
		expected := make([]string, 0, len(s.employeesActive))
		for id := range s.employeesActive {
			expected = append(expected, id)
		}
		all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
		for i := 1; i < len(sorted); i++ {
			if s.employees[sorted[i-1]].CreatedAt.After(s.employees[sorted[i]].CreatedAt) {
				t.Errorf("invalid order by created at asc: %v > %v", s.employees[sorted[i-1]].CreatedAt, s.employees[sorted[i]].CreatedAt)
			}
		}
		for id, e := range all {
			require.Equal(t, s.employees[id], e)
			require.Equal(t, string(serviceModels.EmployeeStatusActive), e.Status)
		}
	})

	t.Run("inactive employees only", func(t *testing.T) {
		limit := 4
		offset := 0
		status := serviceModels.EmployeeStatusInactive
		filter := serviceModels.EmployeeFilters{
			Status: &status,
			BaseList: serviceModels.BaseList{
				Limit:  &limit,
				Offset: &offset,
			},
		}
		expected := make([]string, 0, len(s.employeesIncative))
		for id := range s.employeesIncative {
			expected = append(expected, id)
		}
		all, sorted := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
		for i := 1; i < len(sorted); i++ {
			if s.employees[sorted[i-1]].CreatedAt.After(s.employees[sorted[i]].CreatedAt) {
				t.Errorf("invalid order by created at asc: %v > %v", s.employees[sorted[i-1]].CreatedAt, s.employees[sorted[i]].CreatedAt)
			}
		}
		for id, e := range all {
			require.Equal(t, s.employees[id], e)
			require.Equal(t, string(serviceModels.EmployeeStatusInactive), e.Status)
		}
	})

	for _, role := range []serviceModels.EmployeeRole{serviceModels.EmployeeRoleManager} {
		expected := make([]string, 0)
		for id, e := range s.employees {
			if e.Role == string(role) {
				expected = append(expected, id)
			}
		}
		if len(expected) == 0 {
			t.Errorf("no employees with role %s", role)
			continue
		}
		t.Run("filter by role "+string(role), func(t *testing.T) {
			limit := 4
			offset := 0
			filter := serviceModels.EmployeeFilters{
				Role: &role,
				BaseList: serviceModels.BaseList{
					Limit:  &limit,
					Offset: &offset,
				},
			}
			all, _ := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
			for id, e := range all {
				require.Equal(t, s.employees[id], e)
				require.Equal(t, string(role), e.Role)
			}
		})
	}

	t.Run("filter by full name", func(t *testing.T) {
		// Берем имя первого работника и ищем по части имени
		var searchName string
		for _, e := range s.employees {
			if len(e.FullName) > 3 {
				searchName = e.FullName[:3]
				break
			}
		}
		require.NotEmpty(t, searchName, "no employee with name longer than 3 chars")

		limit := 10
		offset := 0
		filter := serviceModels.EmployeeFilters{
			FullName: &searchName,
			BaseList: serviceModels.BaseList{
				Limit:  &limit,
				Offset: &offset,
			},
		}
		expected := make([]string, 0)
		for id, e := range s.employees {
			if models.MatchesPattern(searchName, e.FullName) {
				expected = append(expected, id)
			}
		}
		all, _ := s.CheckList(t, s.adminToken.AccessToken, &filter, expected)
		for id, e := range all {
			require.Equal(t, s.employees[id], e)
			require.Contains(t, strings.ToLower(e.FullName), strings.ToLower(searchName))
		}
	})
}

func (s *EmployeeSuite) ListInvalid(t *testing.T) {
	t.Run("with invalid order_by", func(t *testing.T) {
		var orderBy serviceModels.EmployeeOrderBy = "invalid"
		filter := serviceModels.EmployeeFilters{OrderBy: &orderBy}
		_, code, err := s.base.client.GetEmployees(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with invalid order_direction", func(t *testing.T) {
		var orderDirection serviceModels.OrderDirection = "invalid"
		filter := serviceModels.EmployeeFilters{
			BaseList: serviceModels.BaseList{
				OrderDirection: &orderDirection,
			},
		}
		_, code, err := s.base.client.GetEmployees(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with invalid role", func(t *testing.T) {
		var role serviceModels.EmployeeRole = "invalid"
		filter := serviceModels.EmployeeFilters{Role: &role}
		_, code, err := s.base.client.GetEmployees(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with invalid status", func(t *testing.T) {
		var status serviceModels.EmployeeStatus = "invalid"
		filter := serviceModels.EmployeeFilters{Status: &status}
		_, code, err := s.base.client.GetEmployees(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with negative limit", func(t *testing.T) {
		limit := -10
		filter := serviceModels.EmployeeFilters{
			BaseList: serviceModels.BaseList{
				Limit: &limit,
			},
		}
		_, code, err := s.base.client.GetEmployees(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("with negative offset", func(t *testing.T) {
		offset := -10
		filter := serviceModels.EmployeeFilters{
			BaseList: serviceModels.BaseList{
				Offset: &offset,
			},
		}
		_, code, err := s.base.client.GetEmployees(s.base.ctx, s.adminToken.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func (s *EmployeeSuite) ListForbidden(t *testing.T) {
	t.Run("list without token", func(t *testing.T) {
		filter := serviceModels.EmployeeFilters{}
		_, code, err := s.base.client.GetEmployees(s.base.ctx, "", &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, code)
	})

	token := s.LoginEmployee(t)

	t.Run("list with manager role", func(t *testing.T) {
		filter := serviceModels.EmployeeFilters{}
		_, code, err := s.base.client.GetEmployees(s.base.ctx, token.AccessToken, &filter)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, code)
	})
}
