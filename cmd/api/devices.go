package main

import (
	"backplate/internal/db"
	"net/http"
)

// swagger:model ListDevicesResponse
type ListDevicesResponse struct {
	Devices []db.Device
}

// swagger:route GET /devices Devices listDevicesHandler
// List all devices
// responses:
//
//	200: ListDevicesResponse
func (app *application) listDevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := app.service.ListDevices(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, ListDevicesResponse{Devices: devices}, make(http.Header))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// swagger:parameters createDeviceHandler
type CreateDeviceRequest struct {
	// in:body
	Body db.CreateDeviceParams
}

// swagger:model CreateDeviceResponse
type CreateDeviceResponse struct {
	Device db.Device
}

// swagger:route POST /devices Devices createDeviceHandler
// Create a new device
// consumes:
// - application/json
// responses:
//
//	201: CreateDeviceResponse
func (app *application) createDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateDeviceRequest
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	device, err := app.service.CreateDevice(r.Context(), input.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, CreateDeviceResponse{Device: device}, make(http.Header))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
