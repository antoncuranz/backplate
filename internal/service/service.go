package service

import (
	"backplate/internal/db"
	"errors"
)

type ServiceConfig struct {
	InboxDir string `envconfig:"INBOX_DIR" default:"inbox"`
	ImageDir string `envconfig:"IMAGES_DIR" default:"images"`
}

type Service struct {
	Store  *db.Queries
	Config ServiceConfig
}

var ErrRecordNotFound = errors.New("record not found")
