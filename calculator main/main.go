package main

import (
	"fmt"

	"github.com/Ladywiz09/pivottechschool/calculator"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(calculator.Add(2, 3))
	fmt.Println(calculator.Subtract(5, 3))
	fmt.Println(calculator.Multiply(2, 3))
	fmt.Println(calculator.Divide(6, 3))
}
