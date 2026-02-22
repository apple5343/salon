package generator

import (
	"salon/internal/models"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit"
)

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

func RandomPassport() models.Passport {
	return models.Passport{
		Series:   strconv.Itoa(gofakeit.Number(1000, 9999)),
		Number:   strconv.Itoa(gofakeit.Number(100000, 999999)),
		IssuedBy: gofakeit.Name(),
	}
}

func (g *Generator) GenerateEmployee() *models.Employee {
	return &models.Employee{
		FullName:     gofakeit.Name(),
		Phone:        RandomPhone(),
		Email:        gofakeit.Email(),
		PasswordHash: gofakeit.Password(true, true, true, true, false, 8),
		Passport:     RandomPassport(),
		Role:         models.EmployeeRoleManager,
		Status:       models.EmployeeStatusActive,
	}
}
