package main

import "backplate/internal/service"

type ServerConfig struct {
	Host string `envconfig:"SERVER_HOST" default:"localhost"`
	Port string `envconfig:"SERVER_PORT" default:"8090"`
}

type DatabaseConfig struct {
	URL      string `envconfig:"DB_URL" default:"127.0.0.1:5432"`
	Database string `envconfig:"DB_DATABASE" default:"backplate"`
	Username string `envconfig:"DB_USERNAME" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"password"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

type config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Service  service.ServiceConfig
}
