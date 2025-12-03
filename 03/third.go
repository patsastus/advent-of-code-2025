package main

import (
	"fmt"
	"bufio"
	"os"
)
func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func part1(scanner *bufio.Scanner) {
	sum := 0
	for scanner.Scan() {
		line := scanner.Text()
		var bestTen, bestOne, index int = -1, -1, -1
		for i, ten := range line[:len(line)-1] {
			if int(ten - '0') > bestTen { 
				bestTen = int(ten - '0')
				index = i 
				if bestTen == 9 {break}
			}
		}
		for _, one := range line[index+1:] {
			if int(one - '0') > bestOne {
				bestOne = int(one - '0')
				if bestOne == 9 {break}
			}
		}
		sum += bestTen * 10 + bestOne
	}
	fmt.Print(sum)
}


func part2(scanner *bufio.Scanner) {
	sum := 0
	const size = 12
	var result []int = make([]int, size)
	for scanner.Scan() {
		line := scanner.Text()
		var index = -1
		for i:=0; i < size; i++ {
			length := size - i - 1		
			var best, tempIndex int = -1, -1
			for j, digit := range line[index+1:len(line)-length] {
				if int(digit - '0') > best { 
					best = int(digit - '0')
					tempIndex = j+index+1 
					if best == 9 {break}
				}
			}
			result[i] = best
			index = tempIndex
		}
		for i:= size-1; i >=0; i-- {
			if result[i] > -1 {
				sum += result[i] * pow(10, size -1 - i)
			}
		}
	}
	fmt.Print(sum)
}

func main() {
	file, err := os.Open("input03.txt")
	if err != nil {panic("")}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//part1(scanner)
	part2(scanner)
}	