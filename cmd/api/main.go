// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-first-application/internal/config"
	"go-first-application/internal/database"
	"go-first-application/internal/handler"
	appMiddleware "go-first-application/internal/middleware"
	repoPostgres "go-first-application/internal/repository/postgres"
	"go-first-application/internal/service"
	"go-first-application/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(cfg.Env)
	log.Info("starting application", "env", cfg.Env, "port", cfg.Port)

	ctx := context.Background()
	dbPool, err := database.NewPostgresPool(ctx, cfg.DB)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	log.Info("connected to postgres database successfully")

	// Repositories
	userRepo := repoPostgres.NewUserRepository(dbPool)
	projectRepo := repoPostgres.NewProjectRepository(dbPool)
	taskRepo := repoPostgres.NewTaskRepository(dbPool)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWT)
	projectService := service.NewProjectService(projectRepo)
	taskService := service.NewTaskService(taskRepo, projectRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)

	// Root Router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Public Routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(appMiddleware.Authenticate(cfg.JWT))

		// Projects
		r.Route("/projects", func(r chi.Router) {
			r.Post("/", projectHandler.Create)
			r.Get("/", projectHandler.List)
			r.Get("/{id}", projectHandler.GetByID)
			r.Put("/{id}", projectHandler.Update)
			r.Delete("/{id}", projectHandler.Delete)

			// Tasks scoped under projects
			r.Post("/{projectId}/tasks", taskHandler.Create)
			r.Get("/{projectId}/tasks", taskHandler.ListByProject)
		})

		// Direct Tasks manipulation
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/{id}", taskHandler.GetByID)
			r.Put("/{id}", taskHandler.Update)
			r.Delete("/{id}", taskHandler.Delete)
		})
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("http server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdownChan
	log.Info("received termination signal, shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("server shutdown complete")
}
