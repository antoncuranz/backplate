package service

import (
	"backplate/internal/db"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (s *Service) CreateDevice(ctx context.Context, params db.CreateDeviceParams) (db.Device, error) {
	return s.Store.CreateDevice(ctx, params)
}

func (s *Service) ListDevices(ctx context.Context) ([]db.Device, error) {
	return s.Store.ListDevices(ctx)
}

func (s *Service) GetDevice(ctx context.Context, id int64) (db.Device, error) {
	device, err := s.Store.GetDevice(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return db.Device{}, ErrRecordNotFound
		default:
			return db.Device{}, err
		}
	}

	return device, nil
}

func (s *Service) UpdateDevice(ctx context.Context, params db.UpdateDeviceParams) (db.Device, error) {
	_, err := s.Store.GetDevice(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return db.Device{}, ErrRecordNotFound
		default:
			return db.Device{}, err
		}
	}

	err = s.Store.UpdateDevice(ctx, params)
	if err != nil {
		return db.Device{}, err
	}

	return s.Store.GetDevice(ctx, params.ID)
}

func (s *Service) DeleteDevice(ctx context.Context, id int64) error {
	_, err := s.Store.GetDevice(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return s.Store.DeleteDevice(ctx, id)
}
