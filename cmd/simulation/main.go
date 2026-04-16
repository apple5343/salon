package main

import (
	"context"
	"log"
	"salon/internal/app"
	"salon/pkg/clock"
	"salon/internal/pkg/simulation"
	"time"

	"go.uber.org/fx"
)

func main() {
	var s *simulation.Simulation
	app := fx.New(
		app.Config(),
		app.Logger(),
		app.InfrastructureModules(),
		app.RepositoryModules(),
		app.ServiceModules(),
		clock.MockClockModule(),
		simulation.SimulationModule(),
		fx.NopLogger,
		fx.Populate(&s),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
	s.Run()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
