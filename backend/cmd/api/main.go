package main

import (
	"context"
	"goapi/internal/api/repository/DAL/SQLite"
	"goapi/internal/api/server"
	"goapi/internal/api/service"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// NewSimpleLogger creates a new log.Logger that writes to a file.
// The file is created if it does not exist, and append to it if it does.
// The file is created with mode 0644.
// The logger also writes to os.Stdout.
// The logger prefixes each log entry with the current date and time and the file name and line number of the calling code.
// The flags argument defines the logging properties.
// The flags are Ldate, Ltime, and Lshortfile.
// The Ldate flag causes the logger to write the current date in the local time zone: 2009/01/23.
// The Ltime flag causes the logger to write the current time in the local time zone: 01:23:23.
// The Lshortfile flag causes the logger to write the file name and line number: logger.go:24.
func NewSimpleLogger(logFile string) *log.Logger {

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)

	}
	return log.New(io.MultiWriter(file, os.Stdout), "", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	// * Timeout is used to gracefully shutdown the server *
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// * Create a logger and database connection *
	logger := NewSimpleLogger("production.log")
	db, err := SQLite.NewSqlite("production.db")
	if err != nil {
		logger.Println("Error setting up database:", err)
		return
	}
	defer db.Close()

	// * Create a service factory and API server *
	sf := service.NewServiceFactory(db, logger, ctx)

	// * Create the API server *
	server := server.NewServer(ctx, sf, logger)

	// * Setup graceful shutdown *
	gracefullShutdown(server, cancel, logger)

	// * Start the server *
	logger.Println("Starting server on :8080...")
	if err := server.ListenAndServe(":8080"); err != nil {
		// If the server was shutdown gracefully, don't log a startup error
		if err != http.ErrServerClosed {
			logger.Println("Server startup error:", err)
		}
		logger.Println("Server gracefully shutdown complete.")
		return
	}
}

func gracefullShutdown(server *server.Server, cancel context.CancelFunc, logger *log.Logger) {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// * Listen for signals to shutdown the server gracefully *
	go func() {
		<-signalCh
		cancel()
		if err := server.Shutdown(); err != nil {
			logger.Println("Error shutting down API Server:", err)
		}
	}()
}
