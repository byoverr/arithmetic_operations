package handlers

import (
	"arithmetic_operations/check_expression"
	"arithmetic_operations/orchestrator/models"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

func HandlerCreateExpression(log *slog.Logger, expressionSaver func(expression *models.Expression) error) http.HandlerFunc {
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

		errValidating := check_expression.CheckExpression(log, inputExpression.Expression)

		if errValidating != nil {
			expression = models.NewExpressionInvalid(inputExpression.Expression)
		} else {
			expression = models.NewExpressionInProcess(inputExpression.Expression)
		}

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

		//TODO: add parser and start to solve

		render.Status(r, http.StatusOK)
		render.JSON(w, r, expression)

		log.Info("expression added", slog.Int("id", expression.Id))
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

		//errValidating := validators.ValidateOperation(operation)

		//if errValidating != nil {
		//	render.Status(r, http.StatusInternalServerError)
		//	render.JSON(w, r, models.NewError(errValidating.Error()))
		//	return
		//}

		errDb := operationUpdate(&operation)

		if errDb != nil {
			log.Error("could not update operation: %+v", operation)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, models.NewError("could not update operation"))
		}

		log.Info("successful to update operation")
		w.WriteHeader(http.StatusOK)
	}
}
