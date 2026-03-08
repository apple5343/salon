package integration

import (
	"net/http"
	"salon/tests/integration/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func (s *BaseTestSuite) LoginAdmin(t *testing.T) *models.EmployeeToken {
	token, code, err := s.client.LoginEmployee(s.ctx, s.adminCreds.Email, s.adminCreds.Password)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)
	return token
}

func (s *BaseTestSuite) HireManager(t *testing.T, adminAccess string) (*models.EmployeeToken, *models.EmployeeCreds) {
	manager := models.GenerateEmployee()
	password := manager.Password
	manager, code, err := s.client.RegisterEmployee(s.ctx, adminAccess, manager)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)
	s.client.HireEmployee(s.ctx, adminAccess, manager.ID)
	token, code, err := s.client.LoginEmployee(s.ctx, manager.Email, password)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)
	return token, &models.EmployeeCreds{Email: manager.Email, Password: password}
}
