// Backplate Api:
//
//	version: 0.0.1
//
// Schemes: http, https
// Host: localhost:8090
// Produces:
//   - application/json
//
// swagger:meta
package main

import (
	"backplate/internal/db"
	"backplate/internal/service"
	"context"
	"errors"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net"
	"net/http"
)

type application struct {
	config  config
	service service.Service
}

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodGet, "/swagger.yaml", http.FileServer(http.Dir("./")))
	router.Handler(http.MethodGet, "/", http.FileServer(http.Dir("./frontend/dist")))

	router.HandlerFunc(http.MethodGet, "/images", app.listImagesHandler)
	router.HandlerFunc(http.MethodPost, "/images", app.createImageHandler)
	router.HandlerFunc(http.MethodGet, "/consume", app.consumeImageHandler)

	router.HandlerFunc(http.MethodPost, "/devices", app.createDeviceHandler)
	router.HandlerFunc(http.MethodGet, "/devices", app.listDevicesHandler)
	router.HandlerFunc(http.MethodGet, "/devices/:id", app.getDeviceHandler)
	router.HandlerFunc(http.MethodPut, "/devices/:id", app.updateDeviceHandler)
	router.HandlerFunc(http.MethodDelete, "/devices/:id", app.deleteDeviceHandler)

	var handler http.Handler = router

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	handler = middleware.SwaggerUI(opts, handler)

	return handler
}

func openDB(config DatabaseConfig) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s",
		config.Username, config.Password, config.URL, config.Database, config.SSLMode)
	println(connString)

	m, err := migrate.New("file://internal/db/migrations", connString)
	if err != nil {
		return nil, err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	err1, err2 := m.Close()
	if err1 != nil || err2 != nil {
		return nil, errors.Join(err1, err2)
	}

	ctx := context.Background()
	return pgx.Connect(ctx, connString)
}

func main() {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	conn, err := openDB(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	app := &application{
		service: service.Service{
			Store:  db.New(conn),
			Config: config.Service,
		},
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Server.Host, config.Server.Port),
		Handler: app.routes(),
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
