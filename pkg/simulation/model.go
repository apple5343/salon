package simulation

import (
	"context"
	"log"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"
)

func (s *Simulation) CreateModel() {
	s.AddEvent(&Event{
		Type: ModelCreated,
		Date: s.currendDay.Format(timeLayout),
	})
}

func (s *Simulation) ProcessModelCreatedEvent(e *Event, t time.Time) {
	adminID, ok := s.RandomAdmin()
	if !ok {
		log.Println("create model: " + ErrNoAvailableAdmins.Error())
		return
	}

	model, err := s.generator.GenerateModel()
	if err != nil {
		log.Println("create model: " + err.Error())
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	_, _, err = s.modelService.Create(ctx, &model)
	if err != nil {
		log.Println("create model: " + err.Error())
		return
	}
}
