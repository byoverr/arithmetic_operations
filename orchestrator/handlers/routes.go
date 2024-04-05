package handlers

import (
	"arithmetic_operations/agent"
	"arithmetic_operations/orchestrator/auth"
	"arithmetic_operations/storage"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func InitRoutes(router *chi.Mux, logger *slog.Logger, repo *storage.PostgresqlDB, agents *agent.Calculator, auth *auth.AuthService) {
	router.Post("/register", HandlerRegisterUser(logger, auth))
	router.Post("/login", HandlerLoginUser(logger, auth))
	router.With(userIdentity).Post("/expression", HandlerCreateExpression(logger, repo.CreateExpression, repo.ReadAllOperations, agents))
	router.With(userIdentity).Get("/expression", HandlerGetAllExpression(logger, repo.ReadAllExpressions))
	router.With(userIdentity).Get("/expression/{id}", HandlerGetExpression(logger, repo.ReadExpression))
	router.With(userIdentity).Get("/operation", HandlerGetAllOperations(logger, repo.ReadAllOperations))
	router.With(userIdentity).Put("/operation", HandlerPutOperations(logger, repo.UpdateOperation))
	router.With(userIdentity).Put("/agent", HandlerAddAgent(logger, agents))
	router.With(userIdentity).Delete("/agent", HandlerRemoveAgent(logger, agents))
	router.With(userIdentity).Get("/agents", HandlerGetAllAgents(logger, agents))
}
