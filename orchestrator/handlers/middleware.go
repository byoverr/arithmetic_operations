package handlers

import (
	"arithmetic_operations/orchestrator/auth"
	"arithmetic_operations/orchestrator/models"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
)

func userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
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
		err := auth.ValidateToken(tokenString, secret)
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
