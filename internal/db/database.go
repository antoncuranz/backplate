package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
)

type DatabaseConfig struct {
	URL      string `envconfig:"DB_URL" default:"127.0.0.1:5432"`
	Database string `envconfig:"DB_DATABASE" default:"backplate"`
	Username string `envconfig:"DB_USERNAME" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"password"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

func ConnectAndMigrate(config DatabaseConfig) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s",
		config.Username, config.Password, config.URL, config.Database, config.SSLMode)

	if err := migrateDatabase(connString); err != nil {
		return nil, err
	}

	return pgx.Connect(context.Background(), connString)
}

func migrateDatabase(connString string) error {
	m, err := migrate.New("file://internal/db/migrations", connString)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	srcErr, dbErr := m.Close()
	return errors.Join(srcErr, dbErr)
}
