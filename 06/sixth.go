package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func part1(inStrings []string) (int, []rune) {
	in := [][]int{}
	operand := []rune{}
	for _, line := range inStrings {
		array := strings.Fields(line)
		if unicode.IsDigit([]rune(array[0])[0]) {
			row := make([]int, len(array))
			for i, v := range array {
				fmt.Sscanf(v, "%d", &row[i])
			}
			in = append(in, row)
		} else {
			operand = make([]rune, len(array))
			for i, v := range array {
				operand[i] = []rune(v)[0]
			}
		}
	}
	sum := 0
	for i := 0; i < len(in[0]); i++ {
		subTotal := in[0][i]
		for j := 1; j < len(in); j++ {
			if operand[i] == '+' {
				subTotal += in[j][i]
			} else if operand[i] == '*' {
				subTotal *= in[j][i]
			}
		}
		sum += subTotal
	}
	fmt.Println("Part 1:", sum)
	return sum, operand
}

func part2(inStrings []string, operands []rune) int {
	splitIndexes := make([]int, len(strings.Fields(inStrings[0])))
	index := 0
	for i := 0; i < len(inStrings[0]); i++ {
		allSpaces := true
		for _, line := range inStrings[:len(inStrings)-1] {
			if []rune(line)[i] != ' ' {
				allSpaces = false
				break
			}
		}
		if allSpaces {
			splitIndexes[index] = i
			index++	
		} 
	}
	splitIndexes[index] = len(inStrings[0])
	index++
	sum := 0
	startIndex := 0
	for i := 0; i < index; i++ {
		numbers := []string{}
		for _, line := range inStrings[:len(inStrings)-1] {
			numbers = append(numbers, line[startIndex:splitIndexes[i]])
		}
		startIndex = splitIndexes[i]
		subTotal := 0
		for k := 0; k < len(numbers[0]); k++ {
			element := 0
			for _, num := range numbers {
				if []rune(num)[k] == ' ' {
					continue
				} else {
					element *= 10
					element += int ([]rune(num)[k] - '0')
				}
			}
			fmt.Println("Element:", element)
			if subTotal == 0 || operands[k] == '+' {
				subTotal += element
			} else if operands[k] == '*' {
				subTotal *= element
			}	
		}
		fmt.Println("SubTotal:", subTotal)
		sum += subTotal
	}
	return sum
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic("")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	inStrings := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		inStrings = append(inStrings, line)
	}
	_, operands := part1(inStrings)
	part2(inStrings, operands)

}