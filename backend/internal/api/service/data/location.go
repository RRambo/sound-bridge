package data

import (
	"context"
	"goapi/internal/api/repository/models"
)

type LocationServiceSQLite struct {
	repo models.LocationRepository
	ctx  context.Context
}

func NewLocationServiceSQLite(repo models.LocationRepository) *LocationServiceSQLite {
	return &LocationServiceSQLite{
		repo: repo,
		ctx:  context.Background(),
	}
}

func (s *LocationServiceSQLite) CreateLocation(location *models.Location) error {
	if location.Name == "" {
		return DataError{Message: "Location name is required"}
	}
	// Set as chosen by default when creating
	location.Chosen = true
	return s.repo.CreateLocation(location, s.ctx)
}

func (s *LocationServiceSQLite) GetAllLocations() ([]*models.Location, error) {
	return s.repo.GetAllLocations(s.ctx)
}

func (s *LocationServiceSQLite) GetChosenLocation() (*models.Location, error) {
	return s.repo.GetChosenLocation(s.ctx)
}

func (s *LocationServiceSQLite) SetChosenLocation(id int) error {
	return s.repo.SetChosenLocation(int64(id), s.ctx)
}

func (s *LocationServiceSQLite) UpdateThreshold(id int, newThreshold float64) error {
	return s.repo.UpdateThreshold(int64(id), newThreshold, s.ctx)
}

func (s *LocationServiceSQLite) DeleteLocation(location *models.Location, ctx context.Context) (int64, error) {
	return s.repo.DeleteLocation(location, ctx)
}
