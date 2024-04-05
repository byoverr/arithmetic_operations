package handlers

import (
	"arithmetic_operations/agent"
	"arithmetic_operations/checker"
	"arithmetic_operations/orchestrator/auth"
	"arithmetic_operations/orchestrator/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

func HandlerCreateExpression(log *slog.Logger, expressionSaver func(expression *models.Expression) error, operationreader func() ([]*models.Operation, error), agents *agent.Calculator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inputExpression models.InputExpression
		var expression *models.Expression

		err := render.DecodeJSON(r.Body, &inputExpression)

		if err != nil {
			jsonError := models.NewError("incorrect JSON file")
			log.Error("incorrect JSON file: %s", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jsonError)

			return
		}

		log.Info("request body decoded")

		errValidating := checker.CheckExpression(log, inputExpression.Expression)

		if errValidating != nil {
			expression = models.NewExpressionInvalid(inputExpression.Expression)
		} else {
			expression = models.NewExpressionInProcess(inputExpression.Expression)
		}
		expression.Expression = checker.RemoveAllSpaces(expression.Expression)
		if errValidating != nil {
			jsonError := models.NewError(errValidating.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jsonError)
			return
		}
		errDb := expressionSaver(expression)

		if errDb != nil {
			log.Error("problem with database", slog.String("error", errDb.Error()))

			jsonError := models.NewError("problem with database")

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, jsonError)

			return
		} else {
			log.Info("added expression to db", expression)
		}

		operations, errDb := operationreader()
		if errDb != nil {
			log.Error("problem with database", slog.String("error", errDb.Error()))

			jsonError := models.NewError("problem with database")

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, jsonError)

			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, expression)

		log.Info("expression added", slog.Int("id", expression.Id))
		agents.CreateTask(expression, operations)
		log.Info("task created")
	}
}

func HandlerGetAllExpression(log *slog.Logger, expressionReader func() ([]*models.Expression, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("start get all expression")

		expressions, err := expressionReader()

		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error to get expression: %s", err)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, models.NewError("no expressions"))
			return
		}

		log.Info("successful to get all expressions")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, expressions)
	}
}

func HandlerGetExpression(log *slog.Logger, expressionReader func(int) (*models.Expression, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		log.Info("start get expression", slog.String("id", idStr))

		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Error("id should be integer and bigger than 0")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, models.NewError("id should be integer"))
			return
		}

		expression, err := expressionReader(id)

		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error to get expression: %s", err)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, models.NewError("no expression with this id"))
			return
		}

		log.Info("successful to get expressions")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, expression)
	}
}

func HandlerGetAllOperations(log *slog.Logger, operationReader func() ([]*models.Operation, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("start get all operations")

		operations, err := operationReader()

		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error to get operations: %s", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, models.NewError("no operations"))
			return
		}

		log.Info("successful to get all operations")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, operations)
	}
}

func HandlerPutOperations(log *slog.Logger, operationUpdate func(operation *models.Operation) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("start put operations")

		var operation models.Operation

		err := render.DecodeJSON(r.Body, &operation)

		if err != nil {
			log.Error("incorrect JSON file: %s", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, models.NewError("incorrect JSON file"))
			return
		}

		log.Info("request body decoded")

		errValidating := checker.ValidateOperation(operation)

		if errValidating != nil {
			log.Error("error with validating operation", operation, slog.String("error", errValidating.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, models.NewError(errValidating.Error()))
			return
		}

		errDb := operationUpdate(&operation)

		if errDb != nil {
			log.Error("could not update operation: ", operation, slog.String("errorDb", errDb.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, models.NewError("could not update operation"))
		}

		log.Info("successful to update operation")
		w.WriteHeader(http.StatusOK)
	}
}

func HandlerAddAgent(log *slog.Logger, calc *agent.Calculator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		calc.AddAgent()
		log.Info("successful to add agent")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, fmt.Sprintf("{'count_of_agents': %d}", len(calc.Agents)))
	}
}

func HandlerRemoveAgent(log *slog.Logger, calc *agent.Calculator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("start removing agent")

		err := calc.RemoveAgent()
		if err != nil {
			log.Error("error with removing: %s", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, models.NewError("error with removing agent (you have only one)"))
			return
		}
		log.Info("successful to remove agent")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, fmt.Sprintf("{'count_of_agents': %d}", calc.NumberOfAgents))
	}
}

func HandlerGetAllAgents(log *slog.Logger, calc *agent.Calculator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("start get all agents")

		agents := calc.Agents
		log.Info("successful to get all agents")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, agents)
	}
}

func HandlerRegisterUser(log *slog.Logger, auth *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inputUser models.RegisterUser

		err := render.DecodeJSON(r.Body, &inputUser)

		if err != nil {
			jsonError := models.NewError("incorrect JSON file")
			log.Error("incorrect JSON file: %s", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jsonError)

			return
		}
		log.Info("request body decoded")

		errValidating := checker.CheckUser(log, &inputUser)
		if errValidating != nil {
			log.Error("error with checking user", slog.String("error", errValidating.Error()))
			jsonError := models.NewError(errValidating.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jsonError)
			return
		}
		id, err := auth.Register(inputUser.Username, inputUser.Password)
		if err != nil {
			log.Error("error with register user", slog.String("error", err.Error()))
			jsonError := models.NewError(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jsonError)
			return
		}
		log.Info("successful to register user")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, map[string]interface{}{"id": id})

	}
}

func HandlerLoginUser(log *slog.Logger, auth *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inputUser models.RegisterUser

		err := render.DecodeJSON(r.Body, &inputUser)

		if err != nil {
			jsonError := models.NewError("incorrect JSON file")
			log.Error("incorrect JSON file: %s", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jsonError)

			return
		}
		log.Debug("request body decoded")

		errValidating := checker.CheckUser(log, &inputUser)
		if errValidating != nil {
			log.Error("error with checking user", slog.String("error", errValidating.Error()))
			jsonError := models.NewError(errValidating.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jsonError)
			return
		}
		token, err := auth.Login(inputUser.Username, inputUser.Password)
		if err != nil {
			log.Error("error with register user", slog.String("error", err.Error()))
			jsonError := models.NewError(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jsonError)
			return
		}
		log.Info("successful to login user")

		w.Header().Add("Authorization", token)
		w.WriteHeader(http.StatusOK)

		render.JSON(w, r, map[string]interface{}{"token": token})

	}

}

//func (h *Handler) userIdentity(c *gin.Context) {
//	header := c.GetHeader(authorizationHeader)
//	if header == "" {
//		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
//		return
//	}
//
//	headerParts := strings.Split(header, " ")
//	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
//		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
//		return
//	}
//
//	if len(headerParts[1]) == 0 {
//		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
//		return
//	}
//
//	userId, err := h.services.Authorization.ParseToken(headerParts[1])
//	if err != nil {
//		newErrorResponse(c, http.StatusUnauthorized, err.Error())
//		return
//	}
//
//	c.Set(userCtx, userId)
//}
