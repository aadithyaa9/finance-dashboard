package main

import (
	"log"
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/aadithyaa9/finance-dashboard/config"
	"github.com/aadithyaa9/finance-dashboard/db"
	"github.com/aadithyaa9/finance-dashboard/internal/auth"
	"github.com/aadithyaa9/finance-dashboard/internal/dashboard"
	"github.com/aadithyaa9/finance-dashboard/internal/middleware"
	"github.com/aadithyaa9/finance-dashboard/internal/records"
	"github.com/aadithyaa9/finance-dashboard/internal/users"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg.DBURL)
	defer database.Close()

	// --- Stores ---
	userStore := users.NewStore(database)
	recordStore := records.NewStore(database)

	// --- Services ---
	authService := auth.NewService(userStore, cfg.JWTSecret, cfg.JWTExpiryHours)
	dashboardService := dashboard.NewService(database)

	// --- Handlers ---
	authHandler := auth.NewHandler(authService)
	recordHandler := records.NewHandler(recordStore)
	userHandler := users.NewHandler(userStore)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	// --- Router ---
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	// Public
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected — all routes below require a valid JWT
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(authService))

		// Records — viewer and above can read; analyst and admin can write; admin can delete
		r.Route("/api/records", func(r chi.Router) {
			r.Get("/", recordHandler.List)
			r.Get("/{id}", recordHandler.GetByID)
			r.With(middleware.RequireRole("analyst", "admin")).Post("/", recordHandler.Create)
			r.With(middleware.RequireRole("analyst", "admin")).Put("/{id}", recordHandler.Update)
			r.With(middleware.RequireRole("admin")).Delete("/{id}", recordHandler.Delete)
		})

		// Dashboard — all authenticated users
		r.Route("/api/dashboard", func(r chi.Router) {
			r.Get("/summary", dashboardHandler.Summary)
			r.Get("/by-category", dashboardHandler.ByCategory)
			r.Get("/trends", dashboardHandler.Trends)
			r.Get("/recent", dashboardHandler.Recent)
		})

		// Users — admin only
		r.Route("/api/users", func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Get("/", userHandler.List)
			r.Patch("/{id}/role", userHandler.UpdateRole)
			r.Patch("/{id}/status", userHandler.UpdateStatus)
		})
	})

	log.Printf("server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
