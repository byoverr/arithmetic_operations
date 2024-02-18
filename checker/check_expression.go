package checker

import (
	"arithmetic_operations/orchestrator/topostfix"
	"arithmetic_operations/stack"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"
)

// RemoveAllSpaces убирает пробелы в выражении
func RemoveAllSpaces(a string) string {
	r := regexp.MustCompile(`\s+`)
	noSpaces := r.ReplaceAllString(a, "")

	return noSpaces
}

// CheckExpression проверяет на все возможные ошибки
//func CheckExpression(log *slog.Logger, expression string) error {
//	var wg sync.WaitGroup
//	log.Info("start check expression", slog.String("expr", expression))
//	if len(expression) == 0 {
//		log.Error("length of expression is 0", slog.String("expr", expression))
//		return errors.New("length of expression is 0")
//	}
//	RemoveAllSpaces(expression)
//	err := make(chan error)
//	ctx := context.Background()
//	wg.Add(7)
//	Checker(&ctx, HasDoubleSymbol, expression, &wg, err)
//	Checker(&ctx, ExpressionStartsWithNumber, expression, &wg, err)
//	Checker(&ctx, IsValidParentheses, expression, &wg, err)
//	Checker(&ctx, HasDivizionByZero, expression, &wg, err)
//	Checker(&ctx, HasValidCharacters, expression, &wg, err)
//	Checker(&ctx, HasAtLeastOneExpression, expression, &wg, err)
//	Checker(&ctx, ContainsCorrectFloatPoint, expression, &wg, err)
//	wg.Wait()
//	if len(err) == 0 {
//		close(err)
//		log.Info("successful check expression", slog.String("expr", expression))
//		return nil
//	} else {
//		errChan := <-err
//		close(err)
//		log.Error("error with checking", slog.String("error", errChan.Error()))
//		return errChan
//	}
//}

func CheckExpression(log *slog.Logger, expression string) error {
	var wg sync.WaitGroup
	log.Info("start check expression", slog.String("expr", expression))
	if len(expression) == 0 {
		log.Error("length of expression is 0", slog.String("expr", expression))
		return errors.New("length of expression is 0")
	}
	RemoveAllSpaces(expression)
	errChan := make(chan error, 7)
	ctx := context.Background()
	wg.Add(7)

	go Checker(&ctx, HasDoubleSymbol, expression, &wg, errChan)
	go Checker(&ctx, ExpressionStartsWithNumber, expression, &wg, errChan)
	go Checker(&ctx, IsValidParentheses, expression, &wg, errChan)
	go Checker(&ctx, HasDivizionByZero, expression, &wg, errChan)
	go Checker(&ctx, HasValidCharacters, expression, &wg, errChan)
	go Checker(&ctx, HasAtLeastOneExpression, expression, &wg, errChan)
	go Checker(&ctx, ContainsCorrectFloatPoint, expression, &wg, errChan)

	wg.Wait()

	if len(errChan) == 0 {
		log.Info("successful check expression", slog.String("expr", expression))
		return nil
	} else {
		err := <-errChan
		log.Error("error with checking", slog.String("error", err.Error()))
		return err
	}
}

type ValidatorFunc func(str string) error

func Checker(ctx *context.Context, check ValidatorFunc, expr string, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	var err error
	go func() {
		select {
		case <-(*ctx).Done():
			return
		default:
			err = check(expr)
			if err != nil {
				errChan <- err
				_, cancel := context.WithCancel(*ctx)
				cancel()
				return
			}
		}
	}()
}

// HasDoubleSymbol проверяет на двойной символ
func HasDoubleSymbol(s string) error {
	var last rune
	for _, r := range s {
		if r == last {
			return errors.New("expression has doubled symbol")
		}
		//if r == 42 || r == 43 || r == 45 || r == 47 { // +-*/
		if r == '+' || r == '-' || r == '*' || r == '/' {
			last = r
		}
	}
	return nil
}

// IsValidParentheses проверяет скобочную последовательность
func IsValidParentheses(s string) error {
	stack := &stack.Stack{}
	countOpen := 0
	countClose := 0

	for _, char := range s {
		if char == '(' {
			countOpen++
		} else if char == ')' {
			countClose++
		}
	}

	if countOpen != countClose {
		return errors.New("expression has invalid parentheses")
	}

	for _, r := range s {
		switch r {
		case '(':
			stack.Push('(')
		case ')':
			if stack.Empty() {
				return errors.New("expression has invalid parentheses")
			}
			stack.Pop()
		default:
			continue // Skip non-parentheses characters
		}
	}

	return nil
}

// HasDivizionByZero проверяет есть ли деление на ноль
func HasDivizionByZero(a string) error {
	if strings.Contains(a, "/0") {
		return errors.New("expression has division by zero")
	}
	return nil
}

// HasValidCharacters проверяет на допустимые символы
func HasValidCharacters(a string) error {
	validChars := "1234567890+-*/(). " // Список допустимых символов
	for _, char := range a {
		//(string(s[j]) == "." && number != "")
		if !strings.ContainsRune(validChars, char) {
			return errors.New(fmt.Sprintf("expression has invalid character %d", char))
		}
	}
	return nil
}

// ContainsCorrectFloatPoint проверяет на точку в правильном месте(должно быть в числе float)
func ContainsCorrectFloatPoint(expr string) error {
	length := len(expr)
	for i := 0; i < length; i++ {
		if string(expr[0]) == "." || string(expr[length-1]) == "." {
			return errors.New("expression has a dot in a wrong place")
		}
		if i > 0 && length-1 > i {
			if string(expr[i]) == "." {
				if !(topostfix.IsOperand(expr[i-1]) && topostfix.IsOperand(expr[i+1])) {
					return errors.New("expression has a dot in a wrong place")
				}
			}
		}
	}

	return nil
}

// HasAtLeastOneExpression проверяет на хотя бы одно выражение число оператор число
func HasAtLeastOneExpression(expr string) error {
	// Regular expression to match the pattern "number operator number"
	pattern := `\d+\s*[\+\-\*/]\s*\d+`
	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	// Check if the expression contains at least one match
	matched := r.MatchString(expr)
	if matched {
		return nil
	} else {
		return errors.New("this string doesn't has at least one expression")
	}

}

// ExpressionStartsWithNumber проверяет - первым ли идёт число или скобка в выражении
func ExpressionStartsWithNumber(expr string) error {
	// Регулярное выражение для проверки начала строки с числом или скобкой
	regexPattern := `^(\()?\d+`
	matched, _ := regexp.MatchString(regexPattern, expr)
	if matched {
		return nil
	} else {
		return errors.New("this string doesn't start with number or bracket")
	}

}

// */+- - двойная штука - ошибка - выполнил
// ( больше чем ) - тоже ошибка - выполнил
// деление на ноль - выполнил
// проверка на символы допустимые - выполнил
// разделение на выражения ((2+2) +2) + 2
