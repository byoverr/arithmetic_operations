package main

import (
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
)

func main() {
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

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	setURLPatterns(router, logger, repo)

	logger.Info("start server", slog.String("address", cfg.HTTPServer.Address))

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

func setURLPatterns(router *chi.Mux, logger *slog.Logger, repo *storage.PostgresqlDB) {
	router.Post("/expression", handlers.HandlerCreateExpression(logger, repo.CreateExpression))
	router.Get("/expression", handlers.HandlerGetAllExpression(logger, repo.ReadAllExpressions))
	router.Get("/expression/{id}", handlers.HandlerGetExpression(logger, repo.ReadExpression))
	router.Get("/operation", handlers.HandlerGetAllOperations(logger, repo.ReadAllOperations))
	router.Put("/operation", handlers.HandlerPutOperations(logger, repo.UpdateOperation))
}

//fmt.Print("Enter infix expression: ")
//infixString, err := topostfix.ReadFromInput()
//
//if err != nil {
//	fmt.Println("Error when scanning input:", err.Error())
//	return
//}
//
//lol := topostfix.ToPostfix(infixString)
//fmt.Println(lol)
//for {
//	k, q := topostfix.GetSubExpressions(lol)
//	t := topostfix.CountSubExpressions(k)
//	lol = topostfix.InsertSubExpressions(t, q)
//	fmt.Println(lol, q)
//	if len(q) == 1 {
//		break
//	}
//}
