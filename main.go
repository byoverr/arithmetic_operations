package main

import (
	"arithmetic_operations/agent"
	"arithmetic_operations/orchestrator/config"
	"arithmetic_operations/orchestrator/handlers"
	"arithmetic_operations/orchestrator/prettylogger"
	"arithmetic_operations/storage"
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

	cfg := config.Load()
	opts := prettylogger.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := prettylogger.NewPrettyHandler(os.Stdout, opts)
	//logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger := slog.New(handler)

	repo, err := storage.PostgresqlOpen(cfg)
	if err != nil {
		log.Fatal(err)
	}
	agents, err := agent.InitializeAgents(cfg.Agent.CountOfAgents)
	if err != nil {
		log.Fatal(err)
	}
	agents.CheckerForNewTasks(repo.UpdateExpression)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	setURLPatterns(router, logger, repo, agents)

	logger.Info("start server", slog.String("address", cfg.HTTPServer.Address))

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
}

func setURLPatterns(router *chi.Mux, logger *slog.Logger, repo *storage.PostgresqlDB, agents *agent.Calculator) {
	router.Post("/expression", handlers.HandlerCreateExpression(logger, repo.CreateExpression,
		repo.ReadAllOperations, agents))
	router.Get("/expression", handlers.HandlerGetAllExpression(logger, repo.ReadAllExpressions))
	router.Get("/expression/{id}", handlers.HandlerGetExpression(logger, repo.ReadExpression))
	router.Get("/operation", handlers.HandlerGetAllOperations(logger, repo.ReadAllOperations))
	router.Put("/operation", handlers.HandlerPutOperations(logger, repo.UpdateOperation))
	router.Put("/agent", handlers.HandlerAddAgent(logger, agents))
	router.Delete("/agent", handlers.HandlerRemoveAgent(logger, agents))
	router.Get("/agents", handlers.HandlerGetAllAgents(logger, agents))
}
