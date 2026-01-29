package app

import "go.uber.org/fx"

func NewApp() *fx.App {
	return fx.New(
		Config(),
		Logger(),
		InfrastructureModules(),
		RepositoryModules(),
		ServiceModules(),
		TransportModules(),
	)
}
