package simulation

import (
	"context"
	"log"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"
)

func (s *Simulation) CreateBrand() {
	s.AddEvent(&Event{
		Type: BrandCreated,
		Date: s.currendDay.Format(timeLayout),
	})
}

func (s *Simulation) ProcessBrandCreatedEvent(e *Event, t time.Time) {
	adminID, ok := s.RandomAdmin()
	if !ok {
		log.Println("create brand: " + ErrNoAvailableAdmins.Error())
	}
	b, err := s.generator.GenerateBrand()
	if err != nil {
		log.Println("create brand: " + err.Error())
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	brand, err := s.brandService.Create(ctx, &b)
	if err != nil {
		log.Println("create brand: " + err.Error())
		return
	}
	err = s.generator.BrandCreated(b.ID, brand)
	if err != nil {
		log.Println("create brand: " + err.Error())
		return
	}
}
