package service

import (
	"backplate/internal/db"
	"context"
)

func (s *Service) CreateDevice(ctx context.Context, params db.CreateDeviceParams) (db.Device, error) {
	return s.Store.CreateDevice(ctx, params)
}

func (s *Service) ListDevices(ctx context.Context) ([]db.Device, error) {
	return s.Store.ListDevices(ctx)
}
