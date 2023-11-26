package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	mutablein "github.com/avearmin/mutable-in"
)

func main() {
	muIn := mutablein.NewMutableIn()
	muIn.Init()
	defer muIn.Close()
	scanner := bufio.NewScanner(muIn)

	carryOver := ""
	for {
		fmt.Print("Calc > ")
		if len(carryOver) != 0 {
			muIn.Write([]byte(carryOver))
		}
		scanner.Scan()
		equation := scanner.Text()
		result := evaluateEquation(equation)
		carryOver = strconv.FormatFloat(result, 'f', -1, 64)
	}
}

func evaluateEquation(equation string) float64 {
	parts := strings.Fields(equation)

	operand1, _ := strconv.ParseFloat(parts[0], 64)
	operator := parts[1]
	operand2, _ := strconv.ParseFloat(parts[2], 64)

	var result float64
	switch operator {
	case "+":
		result = operand1 + operand2
	case "-":
		result = operand1 - operand2
	case "*":
		result = operand1 * operand2
	case "/":
		result = operand1 / operand2
	default:
		fmt.Println("Unsupported operator:", operator)
		os.Exit(1)
	}

	return result
}
