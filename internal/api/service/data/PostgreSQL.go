package data

import (
	"context"
	//"goapi/internal/api/repository/DAL/PostgreSQL"
	"goapi/internal/api/repository/models"
	"time"
)

type DataServicePostgreSQL struct {
	repo         models.DataRepository
	locationRepo models.LocationRepository
}

func NewDataServicePostgreSQL(repo models.DataRepository, locationRepo models.LocationRepository) *DataServicePostgreSQL {
	return &DataServicePostgreSQL{
		repo:         repo,
		locationRepo: locationRepo,
	}
}

func (ds *DataServicePostgreSQL) CleanOldData(ctx context.Context) error {
	// SQL deletes rows where measure_time is older than 6 month
	query := `DELETE FROM data WHERE measure_time < datetime('now', '-6 months');`
	_, err := ds.repo.ExecContext(ctx, query)
	return err
}

func (ds *DataServicePostgreSQL) Create(data *models.Data, ctx context.Context) error {
	// If no room_name provided, set it as current chosen location
	if data.RoomName == "" {
		if ds.locationRepo != nil {
			loc, err := ds.locationRepo.GetChosenLocation(ctx)
			if err == nil && loc != nil && loc.Name != "" {
				data.RoomName = loc.Name
			} else {
				data.RoomName = "Unknown" // fallback if no chosen location
			}
		} else {
			data.RoomName = "Unknown"
		}
	}

	if err := ds.ValidateData(data); err != nil {
		return DataError{Message: "Invalid data: " + err.Error()}
	}
	return ds.repo.Create(data, ctx)
}

func (ds *DataServicePostgreSQL) CreateLatest(data *models.Data, ctx context.Context) error {
	// If no room_name provided, set it as current chosen location
	if data.RoomName == "" {
		if ds.locationRepo != nil {
			loc, err := ds.locationRepo.GetChosenLocation(ctx)
			if err == nil && loc != nil && loc.Name != "" {
				data.RoomName = loc.Name
			} else {
				data.RoomName = "Unknown" // fallback if no chosen location
			}
		} else {
			data.RoomName = "Unknown"
		}
	}

	if err := ds.ValidateData(data); err != nil {
		return DataError{Message: "Invalid data: " + err.Error()}
	}
	return ds.repo.CreateLatest(data, ctx)
}

func (ds *DataServicePostgreSQL) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	data, err := ds.repo.ReadOne(id, ctx)
	if err != nil {
		return nil, err
	}

	_ = data

	// We do something to the data, we deduce something from the data!!!
	// This guides the operation intelligently, for example, if the data is of a certain type, then we do something

	return data, nil
}

func (ds *DataServicePostgreSQL) ReadLatest(id string, ctx context.Context) (*models.Data, error) {
	data, err := ds.repo.ReadLatest(id, ctx)
	if err != nil {
		return nil, err
	}

	_ = data

	return data, nil
}

func (ds *DataServicePostgreSQL) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	return ds.repo.ReadMany(page, rowsPerPage, ctx)
}

func (ds *DataServicePostgreSQL) Update(data *models.Data, ctx context.Context) (int64, error) {
	if err := ds.ValidateData(data); err != nil {
		return 0, DataError{Message: "Invalid data: " + err.Error()}
	}
	return ds.repo.Update(data, ctx)
}

func (ds *DataServicePostgreSQL) Delete(data *models.Data, ctx context.Context) (int64, error) {
	return ds.repo.Delete(data, ctx)
}

func (ds *DataServicePostgreSQL) GetDailySummary(roomName string, date time.Time, ctx context.Context) ([]*models.Data, error) {
	if roomName == "" {
		return nil, DataError{Message: "Room name is required"}
	}

	// Get data for the specified room and date
	return ds.repo.GetDailySummary(roomName, date, ctx)
}

func (ds *DataServicePostgreSQL) GetByRoom(roomName string, ctx context.Context) ([]*models.Data, error) {
	if roomName == "" {
		return nil, DataError{Message: "Room name is required"}
	}

	// Get data for the specified room
	return ds.repo.GetByRoom(roomName, ctx)
}

func (ds *DataServicePostgreSQL) ValidateData(data *models.Data) error {
	var errMsg string
	if data.DeviceID == "" || len(data.DeviceID) > 50 {
		errMsg += "DeviceID is required and must be less than 50 characters. "
	}
	if data.RoomName == "" {
		errMsg += "RoomName is required. "
	}
	if data.SoundLevel < 0 || data.SoundLevel > 150 {
		errMsg += "SoundLevel must be between 0 and 150 dB. "
	}
	if data.Threshold < 0 || data.Threshold > 150 {
		errMsg += "Threshold must be between 0 and 150 dB. "
	}

	if errMsg != "" {
		return DataError{Message: errMsg}
	}

	return nil
}
