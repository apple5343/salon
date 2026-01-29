package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string `env:"HTTP_SERVER_ADDR" env-default:"0.0.0.0:8080"`
}

func HttpServerConfig() (*HttpServer, error) {
	cfg := &HttpServer{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
