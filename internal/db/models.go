// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Device struct {
	ID          int64
	Name        string
	Token       string
	LastSync    pgtype.Timestamp
	SleepsUntil pgtype.Timestamp
}

type Image struct {
	ID            int64
	DeviceID      int64
	Permanent     bool
	DataOriginal  string
	DataProcessed string
}