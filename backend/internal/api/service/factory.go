package service

import (
	"context"
	"goapi/internal/api/repository/DAL"
	"goapi/internal/api/repository/DAL/PostgreSQL"
	"goapi/internal/api/repository/DAL/SQLite"
	service "goapi/internal/api/service/data"
	"log"
	"os"
)

type DataServiceType int

const (
	SQLiteDataService DataServiceType = iota
	PostgreSQLDataService
)

func (d DataServiceType) String() string {
	switch d {
	case SQLiteDataService:
		return "SQLiteDataService"
	case PostgreSQLDataService:
		return "PostgreSQLDataService"
	default:
		return "UnknownDataService"
	}
}

type ServiceFactory struct {
	db     DAL.SQLDatabase
	logger *log.Logger
	ctx    context.Context
}

// * Factory for creating data service *
func NewServiceFactory(db DAL.SQLDatabase, logger *log.Logger, ctx context.Context) *ServiceFactory {
	return &ServiceFactory{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}
}

// CreateDataService returns the appropriate DataService based on the serviceType
func (sf *ServiceFactory) CreateDataService(serviceType DataServiceType) (service.DataService, error) {
	dsType := DataServiceType(serviceType)
	sf.logger.Printf("Creating DataService of type: %s", dsType)
	switch serviceType {
	case SQLiteDataService:
		repo, err := SQLite.NewDataRepository(sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		locationRepo, err := SQLite.NewLocationRepository(sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		ds := service.NewDataServiceSQLite(repo, locationRepo)
		return ds, nil
	case PostgreSQLDataService:
		connStr := os.Getenv("EXTERNAL_DATABASE_URL")
		if connStr == "" {
			sf.logger.Println("Error setting up database: DATABASE_URL environment variable is not set")
			return nil, service.DataError{Message: "EXTERNAL_DATABASE_URL environment variable is not set"}
		}
		repo, err := PostgreSQL.NewDataRepository(connStr, sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		locationRepo, err := PostgreSQL.NewLocationRepository(connStr, sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		// You need to implement NewDataServicePostgreSQL in your service/data package
		ds := service.NewDataServicePostgreSQL(repo, locationRepo)
		return ds, nil
	default:
		return nil, service.DataError{Message: "Invalid data service type."}
	}
}

// Optionally, add a similar CreateLocationService for PostgreSQL if needed
func (sf *ServiceFactory) CreateLocationService(serviceType DataServiceType) (service.LocationService, error) {
	switch serviceType {
	case SQLiteDataService:
		repo, err := SQLite.NewLocationRepository(sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		return service.NewLocationServiceSQLite(repo), nil
	case PostgreSQLDataService:
		connStr := os.Getenv("EXTERNAL_DATABASE_URL")
		if connStr == "" {
			sf.logger.Println("Error setting up database: DATABASE_URL environment variable is not set")
			return nil, service.DataError{Message: "EXTERNAL_DATABASE_URL environment variable is not set"}
		}
		repo, err := PostgreSQL.NewLocationRepository(connStr, sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		// You need to implement NewLocationServicePostgreSQL in your service/data package
		return service.NewLocationServicePostgreSQL(repo), nil
	default:
		return nil, service.DataError{Message: "Invalid location service type."}
	}
}
