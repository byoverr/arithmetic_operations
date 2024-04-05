package checker

import (
	"arithmetic_operations/orchestrator/models"
	"fmt"
	"log/slog"
	"regexp"
	"unicode"
)

func isValidPassword(s string) error {
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
		"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range s {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("password must have at least one %s character", name)
	}
	return nil
}

func isValidUsername(username string) error {
	// Используем регулярное выражение для проверки имени пользователя.
	// Регулярное выражение допускает только буквы (в верхнем и нижнем регистре),
	// цифры, нижнее подчеркивание, @ и точку. Длина имени должна быть от 8 до 40 символов.
	regex := regexp.MustCompile(`^[a-zA-Z0-9_.@]{8,40}$`)
	if !regex.MatchString(username) {
		return fmt.Errorf("username is not valid")
	}
	return nil
}

func CheckUser(log *slog.Logger, user *models.RegisterUser) error {
	log.Info("start check user", slog.String("username", user.Username))

	errUsername := isValidUsername(user.Username)
	if errUsername != nil {
		return fmt.Errorf("error with username: %w", errUsername)
	}
	errPassword := isValidPassword(user.Password)
	if errPassword != nil {
		return fmt.Errorf("error with password: %w", errPassword)
	}
	return nil
}
