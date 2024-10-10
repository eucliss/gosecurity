package alert

import "fmt"

func GreaterInt(first int, second int) bool {
	fmt.Println("GreaterInt", first, second)
	return first > second
}

func LessInt(first int, second int) bool {
	fmt.Println("LessInt", first, second)
	return first < second
}

func EqualsInt(first int, second int) bool {
	fmt.Println("EqualsInt", first, second)
	return first == second
}

func EqualsString(first string, second string) bool {
	fmt.Println("EqualsString", first, second)
	return first == second
}

func NotEqualsString(first string, second string) bool {
	fmt.Println("NotEqualsString", first, second)
	return first != second
}
