package main

import (
	"net/http"
)

// swagger:route GET /consume Consume consumeImageHandler
// Consume an image
func (app *application) consumeImageHandler(w http.ResponseWriter, r *http.Request) {
	imagePath, err := app.service.ConsumeImage()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, imagePath)
}
