package app

import (
	"salon/internal/service/car"
	"salon/internal/service/client"
	"salon/internal/service/employee"
	"salon/internal/service/event"
	"salon/internal/service/sale"
	"salon/internal/service/supplier"

	"go.uber.org/fx"
)

func ServiceModules() fx.Option {
	return fx.Module(
		"service",
		employee.NewModule(),
		supplier.NewModule(),
		client.NewModule(),
		car.NewModule(),
		event.NewModule(),
		sale.NewModule(),
	)
}
