package main

import (
	"errors"
	"fmt"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

func convert(srcPath string, dstPath string) error {
	imageinput, err1 := os.Open(srcPath)
	defer imageinput.Close()

	outfile, err2 := os.Create(dstPath)
	defer outfile.Close()

	src, err3 := jpeg.Decode(imageinput)

	if err := errors.Join(err1, err2, err3); err != nil {
		return err
	}

	resized := resize.Resize(820, 1200, src, resize.Lanczos3)

	gray := image.NewGray(resized.Bounds())
	draw.Draw(gray, gray.Bounds(), resized, resized.Bounds().Min, draw.Src)

	// dither
	grayLevels := 8
	palette := make([]color.Color, grayLevels)
	step := 0xffff / (uint32(grayLevels) - 1)
	for i := 0; i < grayLevels; i++ {
		grayValue := uint16(step * uint32(i))
		grayColor := color.Gray16{Y: grayValue}
		palette[i] = grayColor
	}

	d := dither.NewDitherer(palette)
	d.Matrix = dither.FloydSteinberg
	dithered := d.Dither(gray)

	return bmp.Encode(outfile, dithered)
}

var index = 0

func consume(w http.ResponseWriter, req *http.Request) {
	entries, err1 := os.ReadDir("./images")
	chosen := fmt.Sprintf("./images/%s", entries[index].Name())
	err2 := convert(chosen, "./tmp.bmp")

	if err := errors.Join(err1, err2); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, req, "./tmp.bmp")
	index = (index + 1) % len(entries)
}

func main() {
	http.HandleFunc("/consume", consume)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
