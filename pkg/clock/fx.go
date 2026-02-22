package clock

import "go.uber.org/fx"

func ClockModule() fx.Option {
	return fx.Module(
		"clock",
		fx.Provide(
			fx.Annotate(
				func() *RealClock {
					return &RealClock{}
				},
				fx.As(new(Clock)),
			),
		),
	)
}

func MockClockModule() fx.Option {
	return fx.Module(
		"mock-clock",
		fx.Provide(
			fx.Annotate(
				func() *FakeClock {
					return &FakeClock{}
				},
				fx.As(new(Clock)),
				fx.As(new(MockClock)),
			),
		),
	)
}
