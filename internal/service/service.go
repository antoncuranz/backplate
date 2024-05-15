package service

import "backplate/internal/db"

type ServiceConfig struct {
	InboxDir string `envconfig:"INBOX_DIR" default:"inbox"`
	ImageDir string `envconfig:"IMAGES_DIR" default:"images"`
}
type Service struct {
	Store  *db.Queries
	Config ServiceConfig
}

// Alternatively:
//type Service struct {
//	ImageService ImageService
//}
//
//type ImageService struct {
//	Store db.Queries
//}
