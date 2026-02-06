package models

import (
	serviceModels "salon/internal/models"
	httpModels "salon/internal/transport/http/models"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit"
)

const (
	RoleAdmin = "admin"
)

type EmployeeCreds struct {
	Email    string
	Password string
}

type EmployeeToken struct {
	AccessToken  string
	RefreshToken string
}

func RandomPhone() string {
	var phone strings.Builder
	phone.Grow(11)
	phone.WriteString("89")
	for i := 0; i < 9; i++ {
		n := gofakeit.Number(0, 9)
		phone.WriteString(strconv.Itoa(n))
	}
	return phone.String()
}

func RandomPassport() httpModels.Passport {
	return httpModels.Passport{
		Series:   strconv.Itoa(gofakeit.Number(1000, 9999)),
		Number:   strconv.Itoa(gofakeit.Number(100000, 999999)),
		IssuedBy: gofakeit.Name(),
	}
}

func GenerateEmployee() *httpModels.Employee {
	return &httpModels.Employee{
		FullName: gofakeit.Name(),
		Phone:    RandomPhone(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 8),
		Passport: RandomPassport(),
		Role:     string(serviceModels.EmployeeRoleManager),
		Status:   string(serviceModels.EmployeeStatusActive),
	}
}

func CopyEmployee(e *httpModels.Employee) *httpModels.Employee {
	return &httpModels.Employee{
		ID:       e.ID,
		FullName: e.FullName,
		Phone:    e.Phone,
		Email:    e.Email,
		Password: e.Password,
		Passport: httpModels.Passport{
			Series:   e.Passport.Series,
			Number:   e.Passport.Number,
			IssuedBy: e.Passport.IssuedBy,
		},
		Role:      e.Role,
		Status:    e.Status,
		HireDate:  e.HireDate,
		UpdatedAt: e.UpdatedAt,
		CreatedAt: e.CreatedAt,
	}
}
