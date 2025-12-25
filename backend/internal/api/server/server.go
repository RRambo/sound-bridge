package server

import (
	"context"
	"goapi/internal/api/handlers/data"
	"goapi/internal/api/handlers/locations"
	"goapi/internal/api/middleware"
	"goapi/internal/api/service"
	dataService "goapi/internal/api/service/data"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	ctx        context.Context
	HTTPServer *http.Server
	logger     *log.Logger
}

// NewServer creates a new server instance
func NewServer(ctx context.Context, sf *service.ServiceFactory, logger *log.Logger, serviceType service.DataServiceType) *Server {
	// Create API mux
	apiMux := http.NewServeMux()

	// Create DataService
	ds, err := sf.CreateDataService(serviceType)
	if err != nil {
		logger.Fatalf("Error creating data service: %v", err)
	}

	// Create LocationService
	ls, err := sf.CreateLocationService(serviceType)
	if err != nil {
		logger.Fatalf("Error creating location service: %v", err)
	}

	// Setup handlers
	if err := setupDataHandlers(apiMux, logger, ds); err != nil {
		logger.Fatalf("Error setting up data handlers: %v", err)
	}
	if err := setupLocationHandlers(apiMux, logger, ls); err != nil {
		logger.Fatalf("Error setting up location handlers: %v", err)
	}

	// Schedule daily cleanup of old data (older than 6 months)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				if err := ds.CleanOldData(cleanupCtx); err != nil {
					logger.Println("Error cleaning old data:", err)
				} else {
					logger.Println("Old data cleanup completed")
				}
				cancel()
			case <-ctx.Done():
				return
			}
		}
	}()

	// for serving legacy frontend files
	// Main mux serves frontend static files and mounts API under /api/
	mux := http.NewServeMux()
	frontendDir := filepath.Join("..", "..", "..", "build")
	absFrontendDir, _ := filepath.Abs(frontendDir)
	logger.Println("Serving frontend from:", absFrontendDir)
	//mux.Handle("/", http.FileServer(http.Dir(frontendDir)))

	// Apply authentication & common middleware to API
	middlewares := []middleware.Middleware{
		middleware.BasicAuthenticationMiddleware,
		middleware.CommonMiddleware,
	}
	mux.Handle("/api/", http.StripPrefix("/api", middleware.ChainMiddleware(apiMux, middlewares...)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Normalize and remove leading slash so Join works correctly
		reqPath := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")

		// If root, serve index.html
		if reqPath == "" {
			http.ServeFile(w, r, filepath.Join(absFrontendDir, "index.html"))
			return
		}

		// Candidate file path under the build dir
		filePath := filepath.Join(absFrontendDir, reqPath)

		// If file exists and is not a directory, serve it
		if fi, err := os.Stat(filePath); err == nil && !fi.IsDir() {
			// Optional: cache hashed static assets
			if strings.HasPrefix(reqPath, "static/") || strings.Contains(reqPath, ".") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			http.ServeFile(w, r, filePath)
			return
		}

		// Fallback to index.html for client-side routing
		indexPath := filepath.Join(absFrontendDir, "index.html")
		if _, err := os.Stat(indexPath); err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		http.ServeFile(w, r, indexPath)
	})

	return &Server{
		ctx:    ctx,
		logger: logger,
		HTTPServer: &http.Server{
			Handler: mux,
		},
	}
}

// Shutdown gracefully stops the server
func (api *Server) Shutdown() error {
	api.logger.Println("Gracefully shutting down server...")
	return api.HTTPServer.Shutdown(api.ctx)
}

// ListenAndServe starts the HTTP server
func (api *Server) ListenAndServe(addr string) error {
	api.HTTPServer.Addr = addr
	return api.HTTPServer.ListenAndServe()
}

// ==================== DATA HANDLERS ====================
func setupDataHandlers(mux *http.ServeMux, logger *log.Logger, ds dataService.DataService) error {
	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data.GetHandler(w, r, logger, ds)
		case http.MethodPost:
			data.PostHandler(w, r, logger, ds)
		case http.MethodPut:
			data.PutHandler(w, r, logger, ds)
		case http.MethodOptions:
			data.OptionsHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/data/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data.GetByIDHandler(w, r, logger, ds)
		case http.MethodDelete:
			data.DeleteHandler(w, r, logger, ds)
		case http.MethodOptions:
			data.OptionsHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("GET /data/weekly/{room}", func(w http.ResponseWriter, r *http.Request) {
		data.GetByRoomHandler(w, r, logger, ds)
	})

	mux.HandleFunc("/data/daily/{room}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data.GetDailySummaryHandler(w, r, logger, ds)
		case http.MethodOptions:
			data.OptionsHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return nil
}

// ==================== LOCATION HANDLERS ====================
func setupLocationHandlers(mux *http.ServeMux, logger *log.Logger, ls dataService.LocationService) error {
	mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			locations.GetLocationsHandler(w, r, logger, ls)
		case http.MethodPost:
			locations.CreateLocationHandler(w, r, logger, ls)
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("GET /locations/chosen", func(w http.ResponseWriter, r *http.Request) {
		locations.GetChosenLocationHandler(w, r, logger, ls)
	})

	mux.HandleFunc("/locations/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			locations.UpdateLocationHandler(w, r, logger, ls)
		case http.MethodDelete:
			locations.DeleteHandler(w, r, logger, ls)
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return nil
}
