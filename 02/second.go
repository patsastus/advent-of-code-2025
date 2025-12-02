package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func isDoubled(i int) bool {
	str := fmt.Sprintf("%d", i)
	if len(str) % 2 != 0 { return false }
	firstHalf := str[:len(str)/2]
	secondHalf := str[len(str)/2:]
	//fmt.Printf("first: '%s', second '%s'", firstHalf, secondHalf)
	return strings.EqualFold(firstHalf, secondHalf)
}

func isNRepeat(str string, base int) bool {
	if len(str) % base != 0 { return false }
	splits := len(str)/base
	for i :=0; i< splits; i++ {
		for j := i+1; j < splits; j++ {
			if !strings.EqualFold(str[i*base:(i+1)*base], str[j*base:(j+1)*base]) {		
				return false
			}
		}
	}
	return true
}

func isRepeated(i int) bool {
	str := fmt.Sprintf("%d", i)
	for base :=1; base <= len(str)/2; base++ {
		if isNRepeat(str, base) {
			return true
		}
	}
	return false
}

func main() {
	file, err := os.Open("input02.txt")
	if err != nil {panic("")}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	splits := strings.Split(line, ",")
	sum := 0
	for _, str := range splits {
		var start, end int
		fmt.Sscanf(str, "%d-%d", &start, &end)
		for i := start; i <= end; i++ {
			if isRepeated(i) {
				sum += i
			} else {
			//	fmt.Printf("%d is not doubled\n", i)
			}
		}
	}
	fmt.Println(sum)
}