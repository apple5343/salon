package app

import (
	"salon/internal/service/brand"
	"salon/internal/service/car"
	"salon/internal/service/client"
	"salon/internal/service/employee"
	"salon/internal/service/event"
	"salon/internal/service/model"
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
		brand.NewModule(),
		model.NewModule(),
		car.NewModule(),
		event.NewModule(),
		sale.NewModule(),
	)
}
