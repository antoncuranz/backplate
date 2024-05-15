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
	"github.com/go-openapi/runtime/middleware"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	router := http.NewServeMux()

	router.Handle("GET /swagger.yaml", http.FileServer(http.Dir("./")))
	router.Handle("GET /", http.FileServer(http.Dir("./frontend/dist")))

	router.HandleFunc("GET /images", app.listImagesHandler)
	router.HandleFunc("POST /images", app.createImageHandler)
	router.HandleFunc("GET /consume", app.consumeImageHandler)

	router.HandleFunc("POST /devices", app.createDeviceHandler)
	router.HandleFunc("GET /devices", app.listDevicesHandler)
	router.HandleFunc("GET /devices/{id}", app.getDeviceHandler)
	router.HandleFunc("PUT /devices/{id}", app.updateDeviceHandler)
	router.HandleFunc("DELETE /devices/{id}", app.deleteDeviceHandler)

	var handler http.Handler = router

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	handler = middleware.SwaggerUI(opts, handler)

	return handler
}

func main() {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	conn, err := db.ConnectAndMigrate(config.Database)
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
