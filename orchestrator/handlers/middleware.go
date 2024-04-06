package handlers

import (
	"arithmetic_operations/orchestrator/auth"
	"arithmetic_operations/orchestrator/models"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"strings"
)

func userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			slog.Error("invalid auth header")
			jsonError := models.NewError("invalid auth header")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, jsonError)
			return
		}

		if len(headerParts[1]) == 0 {
			slog.Error("request does not contain an access token")
			jsonError := models.NewError("request does not contain an access token")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, jsonError)
			return
		}

		v := viper.New()
		v.AddConfigPath(".")
		v.SetConfigName("config")
		v.SetConfigType("json")

		if err := v.ReadInConfig(); err != nil {
			slog.Error("Error reading config file: %s\n", err)
			return
		}

		secret := v.GetString("auth_service.secret")
		err := auth.ValidateToken(headerParts[1], secret)
		if err != nil {
			slog.Error("error with validating token:", slog.String("error", err.Error()))
			jsonError := models.NewError(err.Error())
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, jsonError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
