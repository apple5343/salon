package main

import (
	"context"
	"log"
	"salon/internal/config"
	"salon/internal/infrastructure/postgres"
	employeeRepo "salon/internal/repository/employee"
	"salon/internal/service/employee"
	"salon/internal/transport/cli"
	"salon/pkg/clock"
	"time"

	"go.uber.org/fx"
)

func main() {
	var cliApp *cli.CLI

	app := fx.New(
		fx.Provide(config.JWTConfig),
		clock.ClockModule(),
		postgres.NewModule(),
		employeeRepo.NewModule(),
		employee.NewModule(),
		cli.NewModule(),
		fx.NopLogger,
		fx.Populate(&cliApp),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	if err := cliApp.Run(); err != nil {
		log.Fatal(err)
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
