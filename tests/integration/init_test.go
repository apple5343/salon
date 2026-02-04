package integration

import (
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

type TestConfig struct {
	PostgresDSN string `env:"POSTGRES_DSN" required:"true"`
	ServerAddr  string `env:"SERVER_ADDR" required:"true"`
}

type BaseTestSuite struct {
	suite.Suite
	config  *TestConfig
	compose compose.ComposeStack
}
