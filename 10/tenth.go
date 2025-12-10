package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Machine struct {
	lights  []bool
	buttons [][]int
}

func main() {
	start := time.Now()
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
	partOne(scanner)
	fmt.Println("Elapsed time:", time.Since(start))
}

func partOne(scanner *bufio.Scanner) {
	machines := make([]Machine, 0)
	for scanner.Scan() {
		m := Machine{}
		line := scanner.Text()
		splits := strings.Split(line, " ")
		m.lights = make([]bool, len(splits[0])-2)
		m.buttons = make([][]int, len(splits)-2)
		for i, r := range []rune(splits[0])[1 : len(splits[0])-1] {
			if r == '#' {
				m.lights[i] = true
			}
		}
		for i, btnStr := range splits[1 : len(splits)-1] {
			btnSplits := strings.Split(btnStr[1:len(btnStr)-1], ",")
			btns := make([]int, len(btnSplits))
			for j, bs := range btnSplits {
				fmt.Sscanf(bs, "%d", &btns[j])
			}
			m.buttons[i] = btns
		}
		machines = append(machines, m)
		fmt.Printf("Parsed machine: %+v\n", m)
	}
	for _, m := range machines {
		for n := 1; n <= len(m.buttons); n++ {
			combos := makeCombinations(n, m)
			fmt.Printf("Total combinations of length %d: %d\n", n, len(combos))
		}
	}
}

func makeCombinations(n int, m Machine) [][]int {
	fmt.Println("Making combinations of length", n)
	var result [][]int
	indexes := make([]int, n)
	numBtns := len(m.buttons)
	i := 0
	for i >= 0 {
		combo := make([]int, n)
		for j := 0; j < n; j++ {
			combo[j] = indexes[j]
		}
		result = append(result, combo)
		fmt.Println(combo)
		i = n - 1 //starting at last increment position, try to increment
		for i >= 0 && indexes[i] == numBtns-1 {
			i--
		}
		if i < 0 {
			break
		} //if we rolled over all positions, i will be -1 and loop ends
		indexes[i]++
		for j := i + 1; j < n; j++ {
			indexes[j] = indexes[i]
		}
	}
	return result
}
