package auth

import (
	"arithmetic_operations/internal/models"
	"arithmetic_operations/internal/storage"
	"fmt"
	"log/slog"
	"time"
)

type Auth interface {
	Login(name string, password string) (string, error)
	Register(name string, password string) (int64, error)
}

type AuthService struct {
	log      *slog.Logger
	repo     *storage.PostgresqlDB
	tokenTTL time.Duration
	secret   string
	cost     int
}

func NewAuthService(log *slog.Logger, repo *storage.PostgresqlDB, tokenTTL time.Duration, secret string, cost int) *AuthService {
	return &AuthService{
		log:      log,
		repo:     repo,
		tokenTTL: tokenTTL,
		secret:   secret,
		cost:     cost,
	}
}

func (s *AuthService) Register(username string, password string) (int64, error) {

	s.log.Debug("start registration user")

	hashed, err := s.GenerateHash(password)
	if err != nil {
		s.log.Error("err hashing password", slog.String("err", err.Error()))
		return 0, fmt.Errorf("error with hashing: %w", err)
	}
	s.log.Debug("password was hashed")

	user := models.NewUser(username, hashed)

	s.log.Debug("start create user in database")

	id, err := s.repo.CreateUser(user)
	if err != nil {
		//if errors.Is(err, helpers.UsernameExistErr) {
		//	log.Error("user name already exists")
		//	return 0, helpers.UsernameExistErr
		//}

		s.log.Error("failed to create user", slog.String("err", err.Error()))
		return 0, fmt.Errorf("error with database: %w", err)
	}

	s.log.Info("created user", slog.Int64("id", id))
	return id, nil
}

func (s *AuthService) Login(username string, password string) (string, error) {

	s.log.Debug("read user from DB")

	user, err := s.repo.GetUser(username)
	if err != nil {
		//if errors.Is(err, helpers.NoRowsErr) {
		//	log.Info("no such user")
		//	return "", helpers.NoRowsErr
		//}

		s.log.Error("error reading user")
		return "", fmt.Errorf("error with database: %w", err)
	}

	s.log.Debug("comparing password and hash")

	samePassword, err := s.Compare(user.HashPassword, password)
	if err != nil {
		s.log.Error("error comparing password")
		return "", fmt.Errorf("error with comparing password: %w", err)
	}

	if !samePassword {
		s.log.Info("passwords don't match")
		return "", fmt.Errorf("passwords don't match")
	}

	token, err := s.GenerateToken(user, s.tokenTTL, s.secret)
	if err != nil {
		s.log.Error("err to generate token", slog.String("err", err.Error()))
		return "", fmt.Errorf("error generate token: %w", err)
	}

	s.log.Info("token generated", slog.String("token", token))
	return token, nil
}
