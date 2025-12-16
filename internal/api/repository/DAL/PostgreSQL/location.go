package PostgreSQL

import (
	"context"
	"database/sql"
	"goapi/internal/api/repository/DAL"
	"goapi/internal/api/repository/models"
)

type LocationRepository struct {
	sqlDB *sql.DB
	ctx   context.Context
}

func NewLocationRepository(connStr string, sqlDB DAL.SQLDatabase, ctx context.Context) (models.LocationRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &LocationRepository{
		sqlDB: db,
		ctx:   ctx,
	}

	// Create locations table
	_, err = repo.sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS locations (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			chosen BOOLEAN NOT NULL DEFAULT FALSE,
			threshold DOUBLE PRECISION NOT NULL DEFAULT 70.0
		);
	`)
	if err != nil {
		return nil, err
	}

	// Create unique index
	_, err = repo.sqlDB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS only_one_chosen_location 
		ON locations(chosen) WHERE chosen = TRUE;
	`)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *LocationRepository) CreateLocation(location *models.Location, ctx context.Context) error {
	tx, err := r.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Unset current chosen location if adding as chosen
	if location.Chosen {
		_, err = tx.ExecContext(ctx, "UPDATE locations SET chosen = FALSE WHERE chosen = TRUE")
		if err != nil {
			return err
		}
	}

	// Insert new location
	res, err := tx.ExecContext(ctx,
		"INSERT INTO locations (name, chosen, threshold) VALUES ($1, $2, $3)",
		location.Name, location.Chosen, location.Threshold)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	location.ID = int64(id)

	return tx.Commit()
}

func (r *LocationRepository) GetAllLocations(ctx context.Context) ([]*models.Location, error) {
	rows, err := r.sqlDB.QueryContext(ctx,
		"SELECT id, name, chosen, threshold FROM locations ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*models.Location
	for rows.Next() {
		var loc models.Location
		var chosen int
		err := rows.Scan(&loc.ID, &loc.Name, &chosen, &loc.Threshold)
		if err != nil {
			return nil, err
		}
		loc.Chosen = chosen == 1
		locations = append(locations, &loc)
	}
	return locations, nil
}

func (r *LocationRepository) GetChosenLocation(ctx context.Context) (*models.Location, error) {
	row := r.sqlDB.QueryRowContext(ctx,
		"SELECT id, name, chosen, threshold FROM locations WHERE chosen = TRUE")

	var loc models.Location
	var chosen int
	err := row.Scan(&loc.ID, &loc.Name, &chosen, &loc.Threshold)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	loc.Chosen = true
	return &loc, nil
}

func (r *LocationRepository) SetChosenLocation(id int64, ctx context.Context) error {
	tx, err := r.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Unset all chosen
	_, err = tx.ExecContext(ctx, "UPDATE locations SET chosen = FALSE WHERE chosen = TRUE")
	if err != nil {
		return err
	}

	// Set new chosen
	_, err = tx.ExecContext(ctx, "UPDATE locations SET chosen = TRUE WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *LocationRepository) UpdateThreshold(id int64, newThreshold float64, ctx context.Context) error {
	tx, err := r.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE locations SET threshold = $1 WHERE id = $2", newThreshold, id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *LocationRepository) DeleteLocation(location *models.Location, ctx context.Context) (int64, error) {
	res, err := r.sqlDB.ExecContext(ctx, "DELETE FROM locations WHERE id = $1", location.ID)
	if err != nil {
		return 0, err
	}

	// get number of affected rows
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
