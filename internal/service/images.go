package service

import (
	"backplate/internal/db"
	"backplate/internal/img"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

func (s *Service) CreateImage(ctx context.Context, file multipart.File) (db.Image, error) {
	timestamp := strings.Fields(time.Now().String())
	filename := fmt.Sprintf("%s/%s_%s.png", s.Config.InboxDir, timestamp[0], timestamp[1])
	dest, err := os.Create(filename)
	if err != nil {
		return db.Image{}, err
	}

	io.Copy(dest, file)

	params := db.CreateImageParams{
		DeviceID:  0,
		Permanent: true,
	}

	return s.Store.CreateImage(ctx, params)
}

func (s *Service) ListImages(ctx context.Context) ([]db.Image, error) {
	return s.Store.ListImages(ctx)
}

func (s *Service) ConsumeImage() (string, error) {
	chosen, err := s.chooseImage()
	if err != nil {
		return "", err
	}

	err = img.Convert(chosen, TmpImage)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(chosen, s.Config.InboxDir) {
		os.Remove(chosen)
	}

	return TmpImage, nil
}

const TmpImage = "./tmp.bmp"

var index = 0

func (s *Service) chooseImage() (string, error) {
	inboxEntries, _ := listFiles(s.Config.InboxDir)

	if len(inboxEntries) > 0 {
		chosen := fmt.Sprintf("%s/%s", s.Config.InboxDir, inboxEntries[0].Name())
		fmt.Println(chosen)
		return chosen, nil
	}

	entries, err := listFiles(s.Config.ImageDir)
	if err != nil {
		return "", err
	}

	chosen := fmt.Sprintf("%s/%s", s.Config.ImageDir, entries[index].Name())
	index = (index + 1) % len(entries)

	fmt.Println(chosen)
	return chosen, nil
}

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
