package simulation

import (
	"context"
	"log"
	"math/rand/v2"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"
)

type HireEmployeeEvent struct {
	EmployeeID string
	AdminID    string
}

func (s *Simulation) CreateEmployee() {
	s.AddEvent(&Event{
		Type: EmployeeCreated,
		Date: s.currendDay.Format(timeLayout),
	})
}

func (s *Simulation) ProcessEmployeeCreatedEvent(e *Event, t time.Time) {
	adminID, ok := s.RandomAdmin()
	if !ok {
		log.Println("create employee: " + ErrNoAvailableAdmins.Error())
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	employee, err := s.employeeService.Register(ctx, s.generator.GenerateEmployee())
	if err != nil {
		log.Println("register employee: " + err.Error())
		return
	}
	s.employees[employee.ID] = employee
	hireTime := t.Add(s.RandomDurationMinutes(15, 120))
	s.AddEvent(&Event{
		Type: EmployeeHired,
		Data: HireEmployeeEvent{
			EmployeeID: employee.ID,
			AdminID:    adminID,
		},
		Date: hireTime.Format(timeLayout),
		Time: &hireTime,
	})
}

func (s *Simulation) ProcessHireEmployeeEvent(e *Event, t time.Time) {
	data, ok := e.Data.(HireEmployeeEvent)
	if !ok {
		log.Println("invalid event data")
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), data.AdminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	employee, err := s.employeeService.Hire(ctx, data.EmployeeID)
	if err != nil {
		log.Println("hire employee: " + err.Error())
		return
	}

	s.activeEmployees = append(s.activeEmployees, employee.ID)
	s.employees[employee.ID] = employee
}

func (s *Simulation) RandomEmployee() (string, bool) {
	if len(s.activeEmployees) == 0 {
		return "", false
	}

	return s.activeEmployees[rand.IntN(len(s.activeEmployees))], true
}

func (s *Simulation) RandomAdmin() (string, bool) {
	if len(s.admins) == 0 {
		return "", false
	}

	return s.activeAdmins[rand.IntN(len(s.activeAdmins))], true
}

func (s *Simulation) CreateAdmin(t time.Time) {
	e := s.generator.GenerateEmployee()
	e.Role = models.EmployeeRoleAdmin
	s.clock.Set(t)
	e, err := s.employeeService.RegisterAdmin(context.TODO(), e)
	if err != nil {
		log.Println("register admin: " + err.Error())
		return
	}

	s.activeAdmins = append(s.activeAdmins, e.ID)
	s.admins[e.ID] = e
}
