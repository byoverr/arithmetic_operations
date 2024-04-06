package auth

import (
	"arithmetic_operations/internal/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type tokenClaims struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(user *models.User, tokenTTL time.Duration, secret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&tokenClaims{
			Id:       user.Id,
			Username: user.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(signedToken, secret string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return fmt.Errorf("couldn't parse claims")
	}
	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		duration := time.Since(time.Unix(time.Now().Unix()-claims.ExpiresAt.Unix(), 0))

		// Extract hours and minutes from the duration
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60
		return fmt.Errorf("token expired by %dh %dmin %dsec", hours, minutes, seconds)
	}

	return nil

}
