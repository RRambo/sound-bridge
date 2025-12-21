package models

// Though this file is named data, it contains only the sound data.

import (
	"context"
	"time"
)

type Data struct {
	ID          int     `json:"id,omitempty"`
	DeviceID    string  `json:"device_id"`             // Arduino device ID
	RoomName    string  `json:"room_name"`             // Name of the working room
	SoundLevel  float64 `json:"sound_level"`           // Level of sound in dB
	Threshold   float64 `json:"threshold"`             // Threshold level in dB
	MeasureTime string  `json:"measure_time"`          // Time of measurement
	IsAlert     bool    `json:"is_alert"`              // Whether the sound level exceeds the threshold
	Description string  `json:"description"`           // Additional information
	IsPeriodic  bool    `json:"is_periodic,omitempty"` // Is the data constantly/periodically measured
}

type DataRepository interface {
	Create(Data *Data, ctx context.Context) error
	CreateLatest(Data *Data, ctx context.Context) error
	ReadOne(id int, ctx context.Context) (*Data, error)
	ReadLatest(id string, ctx context.Context) (*Data, error)
	ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*Data, error)
	Update(data *Data, ctx context.Context) (int64, error)
	Delete(data *Data, ctx context.Context) (int64, error)
	GetDailySummary(roomName string, date time.Time, ctx context.Context) ([]*Data, error) // To retreive daily summary statistics
	GetByRoom(roomName string, ctx context.Context) ([]*Data, error)                       // To retrieve data by room name
	ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error)
}
