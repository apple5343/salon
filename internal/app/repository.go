package app

import (
	"salon/internal/repository/car"
	"salon/internal/repository/client"
	"salon/internal/repository/employee"
	"salon/internal/repository/event"
	"salon/internal/repository/sale"
	"salon/internal/repository/supplier"

	"go.uber.org/fx"
)

func RepositoryModules() fx.Option {
	return fx.Module(
		"repository",
		employee.NewModule(),
		supplier.NewModule(),
		client.NewModule(),
		car.NewModule(),
		event.NewModule(),
		sale.NewModule(),
	)
}
