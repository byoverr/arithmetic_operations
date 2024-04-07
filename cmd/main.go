package main

import (
	"arithmetic_operations/internal/agent"
	"arithmetic_operations/internal/auth"
	"arithmetic_operations/internal/config"
	"arithmetic_operations/internal/handlers"
	"arithmetic_operations/internal/prettylogger"
	"arithmetic_operations/internal/storage"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	var undoneTasksIDs []string
	ctx := context.Background()
	cfg := config.Load()
	opts := prettylogger.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	}
	handler := prettylogger.NewPrettyHandler(os.Stdout, opts)
	//logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger := slog.New(handler)

	repo, err := storage.PostgresqlOpen(cfg, ctx)
	if err != nil {
		log.Fatal(err)
	}
	agents, err := agent.InitializeAgents(cfg.Agent.CountOfAgents)
	if err != nil {
		log.Fatal(err)
	}
	agents.CheckerForNewTasks(repo.UpdateExpression)

	authService := auth.NewAuthService(logger, repo, cfg.AuthService.TokenTTL, cfg.AuthService.Secret, cfg.AuthService.Cost)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	handlers.InitRoutes(router, logger, repo, agents, authService)

	logger.Info("start server", slog.String("address", cfg.HTTPServer.Address))
	// TODO: hide this shown down
	undoneTasks, err := repo.ReadAllExpressionsUndone()
	if err != nil {
		logger.Error("problem with database", slog.String("error", err.Error()))
		log.Fatal(err)
	}

	operations, err := repo.ReadAllOperations()
	if err != nil {
		logger.Error("problem with database", slog.String("error", err.Error()))
		log.Fatal(err)
	}

	for _, task := range undoneTasks {
		agents.CreateTask(task, operations)
		undoneTasksIDs = append(undoneTasksIDs, strconv.Itoa(task.Id))
	}

	if len(undoneTasksIDs) > 0 {
		logger.Info("undone tasks added", slog.String("id", strings.Join(undoneTasksIDs, ", ")))
	} else {
		logger.Info("no undone tasks")
	}
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start")
	}
	// TODO: add pure shutdown
}
