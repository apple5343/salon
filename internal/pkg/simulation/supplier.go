package simulation

import (
	"context"
	"log"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"
)

func (s *Simulation) CreateSupplier() {
	s.AddEvent(&Event{
		Type: SupplierCreated,
		Date: s.currendDay.Format(timeLayout),
	})
}

func (s *Simulation) ProcessSupplierCreatedEvent(e *Event, t time.Time) {
	adminID, ok := s.RandomAdmin()
	if !ok {
		log.Println("create supplier: " + ErrNoAvailableAdmins.Error())
		return
	}
	sup, err := s.generator.GenerateSupplier()
	if err != nil {
		log.Println("create supplier: " + err.Error())
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	supplier, err := s.supplierService.Create(ctx, &sup)
	if err != nil {
		log.Println("create supplier: " + err.Error())
		return
	}
	err = s.generator.SupplierCreated(sup.ID, supplier)
	if err != nil {
		log.Println("create supplier: " + err.Error())
		return
	}
}
