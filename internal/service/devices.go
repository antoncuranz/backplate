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

func (s *Service) GetDevice(ctx context.Context, id int64) (db.Device, error) {
	return s.Store.GetDevice(ctx, id)
}

func (s *Service) UpdateDevice(ctx context.Context, params db.UpdateDeviceParams) (db.Device, error) {
	err := s.Store.UpdateDevice(ctx, params)
	if err != nil {
		return db.Device{}, err
	}

	return s.Store.GetDevice(ctx, params.ID)
}

func (s *Service) DeleteDevice(ctx context.Context, id int64) error {
	return s.Store.DeleteDevice(ctx, id)
}
