package PostgreSQL

import (
	"context"
	"database/sql"
	"goapi/internal/api/repository/DAL"
	"goapi/internal/api/repository/models"
	"time"

	_ "github.com/lib/pq"
	//"log"
)

type DataRepository struct {
	sqlDB *sql.DB
	createStmt,
	upsertLatestStmt *sql.Stmt
	readStmt,
	ReadLatestStmt *sql.Stmt
	readManyStmt,
	updateStmt,
	deleteStmt *sql.Stmt
	ctx context.Context
}

func NewDataRepository(connStr string, sqlDB DAL.SQLDatabase, ctx context.Context) (*DataRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &DataRepository{
		sqlDB: db,
		ctx:   ctx,
	}

	// Create the data table if it doesn't exist
	if _, err := repo.sqlDB.Exec(`CREATE TABLE IF NOT EXISTS data (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		device_id TEXT NOT NULL DEFAULT 'arduino_001',
		room_name TEXT NOT NULL DEFAULT 'unassigned',
		sound_level REAL NOT NULL DEFAULT 0.0,
		threshold REAL NOT NULL DEFAULT 70.0, 
		measure_time TEXT NOT NULL,
		is_alert INTEGER NOT NULL DEFAULT 0,
		description TEXT DEFAULT ''
	);`); err != nil {
		repo.sqlDB.Close()
		return nil, err
	}

	// Create the latest_data table if it doesn't exist
	// for the latest data per device
	if _, err := repo.sqlDB.Exec(`CREATE TABLE  IF NOT EXISTS latest_data (
		device_id TEXT PRIMARY KEY,
		room_name TEXT NOT NULL DEFAULT 'unassigned',
		sound_level DOUBLE PRECISION NOT NULL DEFAULT 0.0,
		threshold DOUBLE PRECISION NOT NULL DEFAULT 70.0, 
		measure_time TEXT NOT NULL,
		is_alert BOOLEAN NOT NULL DEFAULT FALSE,
		description TEXT DEFAULT ''
	);`); err != nil {
		repo.sqlDB.Close()
		return nil, err
	}

	// * Create needed Prepared SQL statements, this is more efficient than running each query individually
	createStmt, err := repo.sqlDB.Prepare(`INSERT INTO data (
	device_id, room_name, sound_level, threshold, measure_time, is_alert, description)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`)

	if err != nil {
		repo.sqlDB.Close() // Close the database connection if statement preparation fails
		return nil, err
	}
	repo.createStmt = createStmt

	upsertLatestStmt, err := repo.sqlDB.Prepare(`INSERT INTO latest_data (
	device_id, room_name, sound_level, threshold, measure_time, is_alert, description)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT(device_id) DO UPDATE SET
		room_name = excluded.room_name,
		sound_level = excluded.sound_level,
		threshold = excluded.threshold,
		measure_time = excluded.measure_time,
		is_alert = excluded.is_alert,
		description = excluded.description`)
	if err != nil {
		repo.sqlDB.Close() // Close the database connection if statement preparation fails
		return nil, err
	}
	repo.upsertLatestStmt = upsertLatestStmt

	// Read single record
	readStmt, err := repo.sqlDB.Prepare(`SELECT 
	id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description
	FROM data WHERE id = $1`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readStmt = readStmt

	// Read latest record
	ReadLatestStmt, err := repo.sqlDB.Prepare(`SELECT 
	device_id, room_name, sound_level, threshold, measure_time, is_alert, description
	FROM latest_data WHERE device_id = $1`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.ReadLatestStmt = ReadLatestStmt

	// Read multiple records with pagination
	readManyStmt, err := repo.sqlDB.Prepare(`SELECT
	id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description
	FROM data LIMIT $1 OFFSET $2`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readManyStmt = readManyStmt

	// Update record
	updateStmt, err := repo.sqlDB.Prepare(`UPDATE data SET 
	device_id = $1, room_name = $2, sound_level = $3, threshold = $4, 
	measure_time = $5, is_alert = $6, description = $7
	WHERE id = $8`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.updateStmt = updateStmt

	// Delete record
	deleteStmt, err := repo.sqlDB.Prepare(`DELETE FROM data WHERE id = $1`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.deleteStmt = deleteStmt

	go Close(ctx, repo)

	return repo, nil

}

func Close(ctx context.Context, r *DataRepository) {

	<-ctx.Done()
	r.createStmt.Close()
	r.upsertLatestStmt.Close()
	r.readStmt.Close()
	r.ReadLatestStmt.Close()
	r.updateStmt.Close()
	r.deleteStmt.Close()
	r.readManyStmt.Close()
	r.sqlDB.Close()
}

func (r *DataRepository) Create(data *models.Data, ctx context.Context) error {

	// Set default values if not provided
	if data.Threshold == 0 {
		data.Threshold = 70.0
	}
	if data.SoundLevel >= data.Threshold {
		data.IsAlert = true
	}
	if data.MeasureTime == "" {
		data.MeasureTime = time.Now().Format(time.RFC3339)
		//Fill with current time if not provided
	}

	// Execute INSERT with correct field order
	res, err := r.createStmt.ExecContext(ctx,
		data.DeviceID,
		data.RoomName,
		data.SoundLevel,
		data.Threshold,
		data.MeasureTime,
		data.IsAlert,
		data.Description)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	data.ID = int(id)
	return nil
}

func (r *DataRepository) CreateLatest(data *models.Data, ctx context.Context) error {
	// 1. upsert latest
	res, err := r.upsertLatestStmt.ExecContext(ctx,
		data.DeviceID,
		data.RoomName,
		data.SoundLevel,
		data.Threshold,
		data.MeasureTime,
		data.IsAlert,
		data.Description)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	data.ID = int(id)
	return nil
}

func (r *DataRepository) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	row := r.readStmt.QueryRowContext(ctx, id)
	var data models.Data

	// Scan with correct field order matching new schema
	err := row.Scan(
		&data.ID,
		&data.DeviceID,
		&data.RoomName,
		&data.SoundLevel,
		&data.Threshold,
		&data.MeasureTime,
		&data.IsAlert,
		&data.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *DataRepository) ReadLatest(id string, ctx context.Context) (*models.Data, error) {
	row := r.ReadLatestStmt.QueryRowContext(ctx, id)
	var data models.Data

	err := row.Scan(
		&data.DeviceID,
		&data.RoomName,
		&data.SoundLevel,
		&data.Threshold,
		&data.MeasureTime,
		&data.IsAlert,
		&data.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *DataRepository) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	if page < 1 {
		return r.ReadAll()
	}

	offset := rowsPerPage * (page - 1)
	rows, err := r.readManyStmt.QueryContext(ctx, rowsPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.Data
	for rows.Next() {
		var d models.Data
		err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.RoomName,
			&d.SoundLevel,
			&d.Threshold,
			&d.MeasureTime,
			&d.IsAlert,
			&d.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

func (r *DataRepository) ReadAll() ([]*models.Data, error) {
	rows, err := r.sqlDB.QueryContext(context.Background(),
		`SELECT id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description 
		FROM data`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.Data
	for rows.Next() {
		var d models.Data
		err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.RoomName,
			&d.SoundLevel,
			&d.Threshold,
			&d.MeasureTime,
			&d.IsAlert,
			&d.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

func (r *DataRepository) Update(data *models.Data, ctx context.Context) (int64, error) {
	// Update measure_time if modifying
	if data.MeasureTime == "" {
		data.MeasureTime = time.Now().Format(time.RFC3339)
	}

	res, err := r.updateStmt.ExecContext(ctx,
		data.DeviceID,
		data.RoomName,
		data.SoundLevel,
		data.Threshold,
		data.MeasureTime,
		data.IsAlert,
		data.Description,
		data.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (r *DataRepository) Delete(data *models.Data, ctx context.Context) (int64, error) {
	res, err := r.deleteStmt.ExecContext(ctx, data.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (r *DataRepository) GetByRoom(roomName string, ctx context.Context) ([]*models.Data, error) {
	// calculate the last 5 weeks
	startTime := time.Now().AddDate(0, 0, -35) // 35 days = 5 weeks

	rows, err := r.sqlDB.QueryContext(ctx, `
		SELECT id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description 
		FROM data 
		WHERE room_name = $1 AND measure_time >= $2`,
		roomName, startTime.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.Data
	for rows.Next() {
		var d models.Data
		err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.RoomName,
			&d.SoundLevel,
			&d.Threshold,
			&d.MeasureTime,
			&d.IsAlert,
			&d.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

func (r *DataRepository) GetDailySummary(roomName string, date time.Time, ctx context.Context) ([]*models.Data, error) {
	// Calculate start and end of the day
	year, month, day := date.Date()
	location := date.Location()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, location)
	endOfDay := startOfDay.Add(24 * time.Hour)

	startOfDayStr := startOfDay.UTC().Format(time.RFC3339)
	endOfDayStr := endOfDay.UTC().Format(time.RFC3339)

	// Add this hourly_average_sound_level to the query it's been implemented
	// AND hourly_average_sound_level IS NOT NULL

	// Query to get data for the specified room and date range
	query := `
	SELECT id, device_id, room_name, sound_level, threshold, measure_time, is_alert, description 
	FROM data
	WHERE room_name = $1
		AND measure_time >= $2
		AND measure_time < $3
	ORDER BY measure_time ASC
	`
	rows, err := r.sqlDB.QueryContext(ctx, query, roomName, startOfDayStr, endOfDayStr)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*models.Data
	for rows.Next() {
		var d models.Data
		err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.RoomName,
			&d.SoundLevel,
			&d.Threshold,
			&d.MeasureTime,
			&d.IsAlert,
			&d.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

// ExecContext executes an arbitrary SQL statement (used for cleanup, etc.)
func (r *DataRepository) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := r.sqlDB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
