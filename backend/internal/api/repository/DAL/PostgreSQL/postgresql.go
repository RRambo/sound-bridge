package PostgreSQL

import (
	"database/sql"
	"goapi/internal/api/repository/DAL"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQL struct {
	sqlDB          *sql.DB
	dataSourceName string
}

func NewPostgreSQL(dataSourceName string) (DAL.SQLDatabase, error) {
	sqlDB, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	PostgreSQL := &PostgreSQL{
		sqlDB:          sqlDB,
		dataSourceName: dataSourceName,
	}

	return PostgreSQL, nil
}

func (s *PostgreSQL) Connection() *sql.DB {
	return s.sqlDB
}

func (s *PostgreSQL) Close() error {
	return s.sqlDB.Close()
}
