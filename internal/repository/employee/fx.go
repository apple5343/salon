package employee

import (
	"salon/internal/repository/employee/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"employee",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
