package models

import "context"

type Location struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Chosen    bool    `json:"chosen"`
	Threshold float64 `json:"threshold"`
}

type LocationRepository interface {
	CreateLocation(location *Location, ctx context.Context) error
	GetAllLocations(ctx context.Context) ([]*Location, error)
	GetChosenLocation(ctx context.Context) (*Location, error)
	SetChosenLocation(id int64, ctx context.Context) error
	UpdateThreshold(id int64, newThreshold float64, ctx context.Context) error
	DeleteLocation(location *Location, ctx context.Context) (int64, error)
}
