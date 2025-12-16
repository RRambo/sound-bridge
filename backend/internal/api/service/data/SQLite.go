package data

import (
	"context"
	"goapi/internal/api/repository/models"
	"time"
)

// * Implementation of DataService for SQLite database *
type DataServiceSQLite struct {
	repo         models.DataRepository
	locationRepo models.LocationRepository // Add locationRepo for accessing locations
}

func NewDataServiceSQLite(repo models.DataRepository, locationRepo models.LocationRepository) *DataServiceSQLite {
	return &DataServiceSQLite{
		repo:         repo,
		locationRepo: locationRepo,
	}
}
func (ds *DataServiceSQLite) CleanOldData(ctx context.Context) error {
	// SQL deletes rows where measure_time is older than 6 months
	query := `DELETE FROM data WHERE measure_time < datetime('now', '-6 months');`
	_, err := ds.repo.ExecContext(ctx, query)
	return err
}
func (ds *DataServiceSQLite) Create(data *models.Data, ctx context.Context) error {
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

func (ds *DataServiceSQLite) CreateLatest(data *models.Data, ctx context.Context) error {
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

func (ds *DataServiceSQLite) ReadOne(id int, ctx context.Context) (*models.Data, error) {

	data, err := ds.repo.ReadOne(id, ctx)
	if err != nil {
		return nil, err
	}

	_ = data

	// We do something to the data, we deduce something from the data!!!
	// This guides the operation intelligently, for example, if the data is of a certain type, then we do something

	return data, nil
}

func (ds *DataServiceSQLite) ReadLatest(id string, ctx context.Context) (*models.Data, error) {

	data, err := ds.repo.ReadLatest(id, ctx)
	if err != nil {
		return nil, err
	}

	_ = data

	return data, nil
}

func (ds *DataServiceSQLite) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	return ds.repo.ReadMany(page, rowsPerPage, ctx)
}

func (ds *DataServiceSQLite) Update(data *models.Data, ctx context.Context) (int64, error) {

	if err := ds.ValidateData(data); err != nil {
		return 0, DataError{Message: "Invalid data: " + err.Error()}
	}
	return ds.repo.Update(data, ctx)
}

func (ds *DataServiceSQLite) Delete(data *models.Data, ctx context.Context) (int64, error) {
	return ds.repo.Delete(data, ctx)
}

func (ds *DataServiceSQLite) ValidateData(data *models.Data) error {
	var errMsg string
	if data.DeviceID == "" || len(data.DeviceID) > 50 {
		errMsg += "DeviceID is required and must be less than 50 characters. "
	}
	if data.RoomName == "" {
		errMsg += "RoomName is required. "
	}
	// Maybe we need to edit the system around RoomName so typos don't messup the data, for robustness and better usability.
	// Maybe by making a predefined list of room names to choose from in the frontend.
	if data.SoundLevel < 0 || data.SoundLevel > 150 {
		errMsg += "SoundLevel must be between 0 and 150 dB. "
	}
	if data.Threshold < 0 || data.Threshold > 150 {
		errMsg += "Threshold must be between 0 and 150 dB. "
	}

	if errMsg != "" {
		return DataError{Message: errMsg}
	}
	// Apparently there was bug here too, it didn't return error message when there was validation error.

	return nil
}

func (ds *DataServiceSQLite) GetDailySummary(roomName string, date time.Time, ctx context.Context) ([]*models.Data, error) {
	if roomName == "" {
		return nil, DataError{Message: "Room name is required"}
	}

	// Get data for the specified room and date
	return ds.repo.GetDailySummary(roomName, date, ctx)
}

func (ds *DataServiceSQLite) GetByRoom(roomName string, ctx context.Context) ([]*models.Data, error) {
	if roomName == "" {
		return nil, DataError{Message: "Room name is required"}
	}

	// Get data for the specified room
	return ds.repo.GetByRoom(roomName, ctx)
}
