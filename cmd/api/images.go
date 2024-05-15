package main

import (
	"backplate/internal/db"
	"mime/multipart"
	"net/http"
)

// swagger:model ListImagesResponse
type ListImagesResponse struct {
	Images []db.Image
}

// swagger:route GET /images Images listImagesHandler
// List all images
// responses:
//
//	200: ListImagesResponse
func (app *application) listImagesHandler(w http.ResponseWriter, r *http.Request) {
	images, err := app.service.ListImages(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, ListImagesResponse{Images: images}, make(http.Header))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// swagger:parameters createImageHandler
type _ struct {
	// in: formData
	// swagger:file
	ImageFile multipart.File `json:"image"`
}

// swagger:model CreateImageResponse
type CreateImageResponse struct {
	Image db.Image
}

// swagger:route POST /images Images createImageHandler
// Upload a new image
// responses:
//
//	201: CreateImageResponse
func (app *application) createImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("image")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	image, err := app.service.CreateImage(r.Context(), file)

	err = app.writeJSON(w, http.StatusCreated, CreateImageResponse{Image: image}, make(http.Header))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
