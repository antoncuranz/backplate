package main

import (
	"backplate/internal/db"
	"backplate/internal/service"
	"errors"
	"net/http"
)

// swagger:parameters deleteDeviceHandler getDeviceHandler
type _ struct {
	// in: path
	ID int64 `json:"id"`
}

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

// swagger:model GetDeviceResponse
type GetDeviceResponse struct {
	Device db.Device
}

// swagger:route GET /devices/{id} Devices getDeviceHandler
// Get a device
// responses:
//
//	200: GetDeviceResponse
func (app *application) getDeviceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	device, err := app.service.GetDevice(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, GetDeviceResponse{Device: device}, make(http.Header))
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
	err := app.readJSON(w, r, &input.Body)
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

// swagger:parameters updateDeviceHandler
type UpdateDeviceRequest struct {
	// in: path
	ID int64 `json:"id"`
	// in:body
	Body db.UpdateDeviceParams
}

// swagger:model CreateDeviceResponse
type UpdateDeviceResponse struct {
	Device db.Device
}

// swagger:route PUT /devices/{id} Devices updateDeviceHandler
// Update a device
// consumes:
// - application/json
// responses:
//
//	201: UpdateDeviceResponse
func (app *application) updateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var input UpdateDeviceRequest
	var err error

	input.ID, err = app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.readJSON(w, r, &input.Body)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// path id has higher priority than id in body
	input.Body.ID = input.ID

	device, err := app.service.UpdateDevice(r.Context(), input.Body)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, UpdateDeviceResponse{Device: device}, make(http.Header))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// swagger:route DELETE /devices/{id} Devices deleteDeviceHandler
// Delete a device
// responses:
//
//	204:
func (app *application) deleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.service.DeleteDevice(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
