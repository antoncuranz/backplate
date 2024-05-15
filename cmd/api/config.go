package main

import (
	"backplate/internal/db"
	"backplate/internal/service"
)

type ServerConfig struct {
	Host string `envconfig:"SERVER_HOST" default:"localhost"`
	Port string `envconfig:"SERVER_PORT" default:"8090"`
}

type config struct {
	Server   ServerConfig
	Database db.DatabaseConfig
	Service  service.ServiceConfig
}
