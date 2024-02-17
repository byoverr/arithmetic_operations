package topostfix

import (
	"arithmetic_operations/orchestrator/models"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func makeID() string {
	b := make([]byte, 16)
	rand.Read(b)
	randomString := hex.EncodeToString(b)
	return randomString
}

func operator(c string) bool {
	return strings.ContainsAny(c, "+ & - & * & /")
}

func operand(c string) bool {
	_, err := strconv.ParseFloat(c, 64)
	return err == nil
}

func removeElements(slice []string, index int, length int) []string {
	return append(slice[:index], slice[index+length:]...)
}

func insertAt(slice []string, index int, value string) []string {
	// Append the value at the end of the slice.
	slice = append(slice, value)

	// Move the elements starting from index one position to the right.
	copy(slice[index+1:], slice[index:])

	// Set the value at the desired index.
	slice[index] = value

	return slice
}

func GetSubExpressions(str string) ([]models.SubExpression, []string) {
	var i = 0
	var oper string
	var subexpressions []models.SubExpression
	expr := strings.Split(str, " ")

	for len(expr)-2 != i && len(expr) != 1 {
		if operand(expr[i]) && operand(expr[i+1]) && operator(expr[i+2]) {
			id := makeID()
			firstElem, _ := strconv.ParseFloat(expr[i], 64)
			secondElem, _ := strconv.ParseFloat(expr[i+1], 64)
			oper = expr[i+2]

			expr = removeElements(expr, i, 3)
			expr = insertAt(expr, i, id)

			subexpressions = append(subexpressions, models.SubExpression{
				Id:          id,
				Index:       i,
				Leftnum:     firstElem,
				Rightnum:    secondElem,
				Operator:    oper,
				ExtraAction: "",
				Value:       0.00,
				IsDone:      false,
			})
			i = 0
		} else if operand(expr[i]) && operator(expr[i+1]) && operand(expr[i+2]) && operator(expr[i+3]) && expr[i+1] == expr[i+3] {
			id := makeID()
			firstElem, _ := strconv.ParseFloat(expr[i], 64)
			secondElem, _ := strconv.ParseFloat(expr[i+2], 64)
			if expr[i+3] == "+" || expr[i+3] == "*" {
				oper = expr[i+3]
			} else if expr[i+3] == "-" {
				oper = "+"
			} else if expr[i+3] == "/" {
				oper = "*"
			}
			extraAction := expr[i+3]
			expr = removeElements(expr, i, 4)
			expr = insertAt(expr, i, id)

			subexpressions = append(subexpressions, models.SubExpression{
				Id:          id,
				Index:       i,
				Leftnum:     firstElem,
				Rightnum:    secondElem,
				Operator:    oper,
				ExtraAction: extraAction,
				Value:       0.00,
				IsDone:      false,
			})
			i = 0
		} else {
			i++
		}
	}
	return subexpressions, expr
}

func InsertSubExpressions(expr []models.SubExpression, sl []string) string {
	for i := range expr {
		if expr[i].IsDone {
			if expr[i].ExtraAction == "" {
				sl[expr[i].Index] = strconv.FormatFloat(expr[i].Value, 'f', 6, 64)
			} else {
				sl[expr[i].Index] = strconv.FormatFloat(expr[i].Value, 'f', 6, 64) + " " + expr[i].ExtraAction
			}

		}
	}
	return strings.Join(sl, " ")
}
func CountSubExpressions(expr models.SubExpression, operations []*models.Operation) (models.SubExpression, error) {
	var timeForAddition, timeForSubtraction, timeForMultiplication, timeForDivision int
	fmt.Println(operations[0], operations[1], operations[2], operations[3])
	for _, i := range operations {
		if i.OperationKind == models.Addition {
			timeForAddition = i.DurationInMilliSecond
		} else if i.OperationKind == models.Subtraction {
			timeForSubtraction = i.DurationInMilliSecond
		} else if i.OperationKind == models.Multiplication {
			timeForMultiplication = i.DurationInMilliSecond
		} else {
			timeForDivision = i.DurationInMilliSecond
		}
	}
	fmt.Println(timeForAddition, timeForSubtraction, timeForMultiplication, timeForDivision)
	if expr.Operator == "+" {
		time.Sleep(time.Duration(timeForAddition) * time.Millisecond)
		expr.Value = expr.Leftnum + expr.Rightnum
		expr.IsDone = true
	} else if expr.Operator == "-" {
		time.Sleep(time.Duration(timeForSubtraction) * time.Millisecond)
		expr.Value = expr.Leftnum - expr.Rightnum
		expr.IsDone = true
	} else if expr.Operator == "*" {
		time.Sleep(time.Duration(timeForMultiplication) * time.Millisecond)
		expr.Value = expr.Leftnum * expr.Rightnum
		expr.IsDone = true
	} else if expr.Operator == "/" {
		time.Sleep(time.Duration(timeForDivision) * time.Millisecond)
		if expr.Rightnum == 0 {
			expr.Value = 0
			expr.IsDone = true
			return models.SubExpression{}, errors.New("division by zero in subexpression")
		}
		expr.Value = expr.Leftnum / expr.Rightnum
		expr.IsDone = true
	}
	return expr, nil
}
