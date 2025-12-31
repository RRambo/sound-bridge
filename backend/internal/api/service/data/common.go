package data

import (
	"context"
	"goapi/internal/api/repository/models"
	"time"
)

type DataService interface {
	Create(data *models.Data, ctx context.Context) error
	CreateLatest(data *models.Data, ctx context.Context) error
	ReadOne(id int, ctx context.Context) (*models.Data, error)
	ReadLatest(id string, ctx context.Context) (*models.Data, error)
	ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error)
	Update(data *models.Data, ctx context.Context) (int64, error)
	Delete(data *models.Data, ctx context.Context) (int64, error)
	ValidateData(data *models.Data) error
	GetDailySummary(roomName string, date time.Time, ctx context.Context) ([]*models.Data, error)
	GetByRoom(roomName string, ctx context.Context) ([]*models.Data, error)
	CleanOldData(ctx context.Context) error
}

type DataError struct {
	Message string
}

func (de DataError) Error() string {
	return de.Message
}
