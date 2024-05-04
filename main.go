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
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const InboxDir = "./inbox"
const ImageDir = "./images"
const TmpImage = "./tmp.bmp"

var index = 0

func listFiles(folder string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var result []os.DirEntry
	for _, file := range files {
		// ignore hidden files (.gitkeep, .DS_store)
		if !strings.HasPrefix(file.Name(), ".") {
			result = append(result, file)
		}
	}

	return result, nil
}

func chooseImage() (string, error) {
	inboxEntries, _ := listFiles(InboxDir)

	if len(inboxEntries) > 0 {
		chosen := fmt.Sprintf("%s/%s", InboxDir, inboxEntries[0].Name())
		fmt.Println(chosen)
		return chosen, nil
	}

	entries, err := listFiles(ImageDir)
	if err != nil {
		return "", err
	}

	chosen := fmt.Sprintf("%s/%s", ImageDir, entries[index].Name())
	index = (index + 1) % len(entries)

	fmt.Println(chosen)
	return chosen, nil
}

func convert(srcPath string, dstPath string) error {
	imageinput, err1 := os.Open(srcPath)
	defer imageinput.Close()

	outfile, err2 := os.Create(dstPath)
	defer outfile.Close()

	src, _, err3 := image.Decode(imageinput)

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
	palette2 := []color.Color{
		color.Gray16{Y: 0x0000},
		color.Gray16{Y: 0x4b4b},
		color.Gray16{Y: 0x5555},
		color.Gray16{Y: 0x6b6b},
		color.Gray16{Y: 0x8282},
		color.Gray16{Y: 0x9696},
		color.Gray16{Y: 0xbaba},
		color.Gray16{Y: 0xd2d2},
	}

	d := dither.NewDitherer(palette2)
	d.Matrix = dither.FloydSteinberg
	dithered := d.DitherPaletted(gray)
	dithered.Palette = palette

	return bmp.Encode(outfile, dithered)
}

func consume(w http.ResponseWriter, req *http.Request) {
	chosen, err1 := chooseImage()
	err2 := convert(chosen, TmpImage)
	if strings.HasPrefix(chosen, InboxDir) {
		os.Remove(chosen)
	}

	if err := errors.Join(err1, err2); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, req, TmpImage)
}

func upload(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	req.ParseMultipartForm(10 << 20)

	file, _, err := req.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer file.Close()

	timestamp := strings.Fields(time.Now().String())
	filename := fmt.Sprintf("%s/%s_%s.png", InboxDir, timestamp[0], timestamp[1])
	dest, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.Copy(dest, file)
}

func main() {
	frontend := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", frontend)

	http.HandleFunc("/consume", consume)
	http.HandleFunc("/upload", upload)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
