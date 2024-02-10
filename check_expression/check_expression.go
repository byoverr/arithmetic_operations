package check_expression

import (
	"arithmetic_operations/stack"
	"regexp"
	"strings"
)

// RemoveAllSpaces убирает пробелы в выражении
func RemoveAllSpaces(a string) string {
	r := regexp.MustCompile(`\s+`)
	noSpaces := r.ReplaceAllString(a, "")

	return noSpaces
}

// HasDoubleSymbol проверяет на двойной символ
func HasDoubleSymbol(s string) bool {
	var last rune
	for _, r := range s {
		if r == last {
			return true
		}
		//if r == 42 || r == 43 || r == 45 || r == 47 { // +-*/
		if r == '+' || r == '-' || r == '*' || r == '/' {
			last = r
		}
	}
	return false
}

// IsValidParentheses проверяет скобочную последовательность
func IsValidParentheses(s string) bool {
	stack := &stack.Stack{}

	for _, r := range s {
		switch r {
		case '(':
			stack.Push('(')
		case ')':
			if stack.Empty() {
				return false
			}
			stack.Pop()
		default:
			continue // Skip non-parentheses characters
		}
	}

	return stack.Empty()
}

// HasDivizionByZero проверяет есть ли деление на ноль
func HasDivizionByZero(a string) bool {
	if strings.Contains(a, "/0") {
		return true
	}
	return false
}

// HasValidCharacters проверяет на допустимые символы
func HasValidCharacters(a string) (bool, rune) {
	validChars := "1234567890+-*/() " // Список допустимых символов
	for _, char := range a {
		if !strings.ContainsRune(validChars, char) {
			return false, char
		}
	}
	return true, 0
}

// */+- - двойная штука - ошибка - выполнил
// ( больше чем ) - тоже ошибка - выполнил
// деление на ноль - выполнил
// проверка на символы допустимые - выполнил
// разделение на выражения ((2+2) +2) + 2
