package data

import (
	"context"
	"goapi/internal/api/repository/models"
	"time"
)

// ================= MOCK SUCCESS =================
type MockDataServiceSuccessful struct{}

func (m *MockDataServiceSuccessful) Create(d *models.Data, ctx context.Context) error {
	return nil
}

func (m *MockDataServiceSuccessful) CreateLatest(d *models.Data, ctx context.Context) error {
	return nil
}
func (m *MockDataServiceSuccessful) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	return &models.Data{
		ID:          id,
		DeviceID:    "arduino_mock",
		RoomName:    "Room_A",
		SoundLevel:  60.0,
		Threshold:   70.0,
		MeasureTime: time.Now().Format(time.RFC3339),
		IsAlert:     false,
		Description: "mock data",
	}, nil
}

func (m *MockDataServiceSuccessful) ReadLatest(id string, ctx context.Context) (*models.Data, error) {
	if id == "invalid" {
		return nil, &DataError{Message: "invalid id"}
	} else {
		return &models.Data{
			DeviceID:    id,
			RoomName:    "Room_A",
			SoundLevel:  60.0,
			Threshold:   70.0,
			MeasureTime: time.Now().Format(time.RFC3339),
			IsAlert:     false,
			Description: "mock data",
		}, nil
	}
}

func (m *MockDataServiceSuccessful) ReadMany(page, rows int, ctx context.Context) ([]*models.Data, error) {
	return []*models.Data{
		{
			ID:          1,
			DeviceID:    "arduino_mock",
			RoomName:    "Room_A",
			SoundLevel:  60.0,
			Threshold:   70.0,
			MeasureTime: time.Now().Format(time.RFC3339),
			IsAlert:     false,
			Description: "mock data",
		},
		{
			ID:          2,
			DeviceID:    "arduino_mock",
			RoomName:    "Room_B",
			SoundLevel:  55.5,
			Threshold:   70.0,
			MeasureTime: time.Now().Format(time.RFC3339),
			IsAlert:     false,
			Description: "mock data 2",
		},
	}, nil
}
func (m *MockDataServiceSuccessful) Update(d *models.Data, ctx context.Context) (int64, error) {
	return 1, nil
}
func (m *MockDataServiceSuccessful) Delete(d *models.Data, ctx context.Context) (int64, error) {
	return 1, nil
}
func (m *MockDataServiceSuccessful) ValidateData(d *models.Data) error {
	return nil
}
func (m *MockDataServiceSuccessful) GetDailySummary(room string, date time.Time, ctx context.Context) ([]*models.Data, error) {
	return []*models.Data{
		{
			ID:          1,
			DeviceID:    "arduino_mock",
			RoomName:    room,
			SoundLevel:  60.0,
			Threshold:   70.0,
			MeasureTime: date.Format(time.RFC3339),
			IsAlert:     false,
			Description: "daily summary mock",
		},
	}, nil
}
func (m *MockDataServiceSuccessful) GetByRoom(room string, ctx context.Context) ([]*models.Data, error) {
	return []*models.Data{
		{
			ID:          1,
			DeviceID:    "arduino_mock",
			RoomName:    room,
			SoundLevel:  60.0,
			Threshold:   70.0,
			MeasureTime: time.Now().Format(time.RFC3339),
			IsAlert:     false,
			Description: "weekly summary mock",
		},
	}, nil
}
func (m *MockDataServiceSuccessful) CleanOldData(ctx context.Context) error {
	return nil
}

// ================= MOCK ERROR =================
type MockDataServiceError struct{}

func (m *MockDataServiceError) Create(d *models.Data, ctx context.Context) error {
	return &DataError{Message: "Error creating data."}
}
func (m *MockDataServiceError) CreateLatest(d *models.Data, ctx context.Context) error {
	return &DataError{Message: "Error creating data."}
}
func (m *MockDataServiceError) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	return nil, &DataError{Message: "Error reading data."}
}
func (m *MockDataServiceError) ReadLatest(id string, ctx context.Context) (*models.Data, error) {
	return nil, &DataError{Message: "Error reading data."}
}
func (m *MockDataServiceError) ReadMany(page, rows int, ctx context.Context) ([]*models.Data, error) {
	return nil, &DataError{Message: "Error reading data."}
}
func (m *MockDataServiceError) Update(d *models.Data, ctx context.Context) (int64, error) {
	return 0, &DataError{Message: "Error updating data."}
}
func (m *MockDataServiceError) Delete(d *models.Data, ctx context.Context) (int64, error) {
	return 0, &DataError{Message: "Internal Server error"}
}
func (m *MockDataServiceError) ValidateData(d *models.Data) error {
	return nil
}
func (m *MockDataServiceError) GetDailySummary(room string, date time.Time, ctx context.Context) ([]*models.Data, error) {
	return nil, &DataError{Message: "Error fetching daily summary."}
}
func (m *MockDataServiceError) GetByRoom(room string, ctx context.Context) ([]*models.Data, error) {
	return nil, &DataError{Message: "Error fetching weekly summary."}
}
func (m *MockDataServiceError) CleanOldData(ctx context.Context) error {
	return &DataError{Message: "Error cleaning old data."}
}

// ================= MOCK NOT FOUND =================
type MockDataServiceNotFound struct{}

func (m *MockDataServiceNotFound) Create(d *models.Data, ctx context.Context) error {
	return nil
}
func (m *MockDataServiceNotFound) CreateLatest(d *models.Data, ctx context.Context) error {
	return nil
}
func (m *MockDataServiceNotFound) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	return nil, nil
}
func (m *MockDataServiceNotFound) ReadLatest(id string, ctx context.Context) (*models.Data, error) {
	return nil, nil
}
func (m *MockDataServiceNotFound) ReadMany(page, rows int, ctx context.Context) ([]*models.Data, error) {
	return []*models.Data{}, nil
}
func (m *MockDataServiceNotFound) Update(d *models.Data, ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *MockDataServiceNotFound) Delete(d *models.Data, ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *MockDataServiceNotFound) ValidateData(d *models.Data) error {
	return nil
}
func (m *MockDataServiceNotFound) GetDailySummary(room string, date time.Time, ctx context.Context) ([]*models.Data, error) {
	return nil, nil
}
func (m *MockDataServiceNotFound) GetByRoom(room string, ctx context.Context) ([]*models.Data, error) {
	return nil, nil
}
func (m *MockDataServiceNotFound) CleanOldData(ctx context.Context) error {
	return &DataError{Message: "Resource not found."}
}
