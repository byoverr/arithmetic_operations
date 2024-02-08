package main

import (
	"fmt"
	"regexp"
)

func removeAllSpaces(a string) string {
	r := regexp.MustCompile(`\s+`)
	noSpaces := r.ReplaceAllString(a, "")

	return noSpaces
}
func hasDoubleSymbol(s string) bool {
	var last rune
	for _, r := range s {
		if r == last {
			return true
		}
		if r == 42 || r == 43 || r == 45 || r == 47 { // +-*/
			last = r
		}
	}
	return false
}

// */+- - двойная штука - ошибка - выполнил
// ( больше чем ) - тоже ошибка
// проверка на символы допустимые
// разделение на выражения ((2+2) +2) + 2
func main() {

	// Исходная строка с пробелами
	original := "2 (( 2 " // ** - двойная штука - ошибка ( больше чем ) - тоже ошибка и + проверка на символы допустимые 43 45 42 47
	fmt.Println(hasDoubleSymbol(original))
}
