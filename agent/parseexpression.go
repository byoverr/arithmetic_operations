package main

import (
	"arithmetic_operations/stack"
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SubExpression struct {
	id          string
	index       int
	leftnum     float64
	rightnum    float64
	operator    string
	extraAction string
	value       float64
	isDone      bool
}

func MakeID() string {
	b := make([]byte, 16)
	rand.Read(b)
	randomString := hex.EncodeToString(b)
	return randomString
}
func Operator(c string) bool {
	return strings.ContainsAny(c, "+ & - & * & /")
}

func Operand(c string) bool {
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

func GetSubExpressions(str string) ([]SubExpression, []string) {
	var i = 0
	var operator string
	var subexpressions []SubExpression
	expr := strings.Split(str, " ")

	for len(expr)-2 != i && len(expr) != 1 {
		if Operand(expr[i]) && Operand(expr[i+1]) && Operator(expr[i+2]) {
			id := MakeID()
			firstElem, _ := strconv.ParseFloat(expr[i], 64)
			secondElem, _ := strconv.ParseFloat(expr[i+1], 64)
			operator = expr[i+2]

			expr = removeElements(expr, i, 3)
			expr = insertAt(expr, i, id)

			subexpressions = append(subexpressions, SubExpression{
				id:          id,
				index:       i,
				leftnum:     firstElem,
				rightnum:    secondElem,
				operator:    operator,
				extraAction: "",
				value:       0.00,
				isDone:      false,
			})
			i = 0
		} else if Operand(expr[i]) && Operator(expr[i+1]) && Operand(expr[i+2]) && Operator(expr[i+3]) && expr[i+1] == expr[i+3] {
			id := MakeID()
			firstElem, _ := strconv.ParseFloat(expr[i], 64)
			secondElem, _ := strconv.ParseFloat(expr[i+2], 64)
			if expr[i+3] == "+" || expr[i+3] == "*" {
				operator = expr[i+3]
			} else if expr[i+3] == "-" {
				operator = "+"
			} else if expr[i+3] == "/" {
				operator = "*"
			}
			extraAction := expr[i+3]
			expr = removeElements(expr, i, 4)
			expr = insertAt(expr, i, id)

			subexpressions = append(subexpressions, SubExpression{
				id:          id,
				index:       i,
				leftnum:     firstElem,
				rightnum:    secondElem,
				operator:    operator,
				extraAction: extraAction,
				value:       0.00,
				isDone:      false,
			})
			i = 0
		} else {
			i++
		}
	}
	return subexpressions, expr
}

func InsertSubExpressions(expr []SubExpression, sl []string) string {
	for i := range expr {
		if expr[i].isDone {
			if expr[i].extraAction == "" {
				sl[expr[i].index] = strconv.FormatFloat(expr[i].value, 'f', 6, 64)
			} else {
				sl[expr[i].index] = strconv.FormatFloat(expr[i].value, 'f', 6, 64) + " " + expr[i].extraAction
			}

		}
	}
	return strings.Join(sl, " ")
}
func CountSubExpressions(expr []SubExpression) []SubExpression {
	for i := range expr {
		if expr[i].operator == "+" {
			expr[i].value = expr[i].leftnum + expr[i].rightnum
			expr[i].isDone = true
		} else if expr[i].operator == "-" {
			expr[i].value = expr[i].leftnum - expr[i].rightnum
			expr[i].isDone = true
		} else if expr[i].operator == "*" {
			expr[i].value = expr[i].leftnum * expr[i].rightnum
			expr[i].isDone = true
		} else if expr[i].operator == "/" {
			expr[i].value = expr[i].leftnum / expr[i].rightnum
			expr[i].isDone = true
		} else {
			continue
		}
	}
	return expr
}

func IsOperator(c uint8) bool {
	return strings.ContainsAny(string(c), "+ & - & * & /")
}

func IsOperand(c uint8) bool {
	return c >= '0' && c <= '9'
}

func getOperatorWeight(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return -1
}

func hasHigherPrecedence(op1 string, op2 string) bool {
	op1Weight := getOperatorWeight(op1)
	op2Weight := getOperatorWeight(op2)
	return op1Weight >= op2Weight
}

func ToPostfix(s string) string {

	var stack stack.Stack

	postfix := ""

	length := len(s)

	for i := 0; i < length; i++ {

		char := string(s[i])
		//// Skip whitespaces
		//if char == " " {
		//	continue
		//}

		if char == "(" {
			stack.Push(char)
		} else if char == ")" {
			for !stack.Empty() {
				str, _ := stack.Top().(string)
				if str == "(" {
					break
				}
				postfix += " " + str
				stack.Pop()
			}
			stack.Pop()
		} else if !IsOperator(s[i]) {
			// If character is not an operator
			// Keep in mind it's just an operand
			j := i
			number := ""
			for ; j < length && IsOperand(s[j]); j++ {
				number = number + string(s[j])
			}
			postfix += " " + number
			i = j - 1
		} else {
			// If character is operator, pop two elements from stack,
			// perform operation and push the result back.
			for !stack.Empty() {
				top, _ := stack.Top().(string)
				if top == "(" || !hasHigherPrecedence(top, char) {
					break
				}
				postfix += " " + top
				stack.Pop()
			}
			stack.Push(char)
		}
	}

	for !stack.Empty() {
		str, _ := stack.Pop().(string)
		postfix += " " + str
	}

	return strings.TrimSpace(postfix)
}

func ReadFromInput() (string, error) {

	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')

	return strings.TrimSpace(s), err
}

func main() {

	fmt.Print("Enter infix expression: ")
	infixString, err := ReadFromInput()

	if err != nil {
		fmt.Println("Error when scanning input:", err.Error())
		return
	}

	lol := ToPostfix(infixString)
	for {
		k, q := GetSubExpressions(lol)
		t := CountSubExpressions(k)
		lol = InsertSubExpressions(t, q)
		fmt.Println(lol, q)
		if len(q) == 1 {
			break
		}
	}
}
