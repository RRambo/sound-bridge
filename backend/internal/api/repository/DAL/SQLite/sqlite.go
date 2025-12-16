package SQLite

import (
	"database/sql"
	"goapi/internal/api/repository/DAL"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	sqlDB          *sql.DB
	dataSourceName string
}

func NewSqlite(dataSourceName string) (DAL.SQLDatabase, error) {

	sqlDB, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	Sqlite := &SQLite{
		sqlDB:          sqlDB,
		dataSourceName: dataSourceName,
	}

	return Sqlite, nil
}

func (s *SQLite) Connection() *sql.DB {
	return s.sqlDB
}

func (s *SQLite) Close() error {
	return s.sqlDB.Close()
}
