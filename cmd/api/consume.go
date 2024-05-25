package main

import (
	"net/http"
)

// swagger:parameters consumeImageHandler
type _ struct {
	// in:query
	Token string `json:"token"`
}

// swagger:route GET /consume Consume consumeImageHandler
// Consume an image
func (app *application) consumeImageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	imagePath, err := app.service.ConsumeImage(token)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Inkplate-Sleep-For-Millis", "1234")
	w.Header().Set("Inkplate-Should-Update-Display", "true")
	http.ServeFile(w, r, imagePath)
}
