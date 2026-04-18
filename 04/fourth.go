package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var filename string
	if len(os.Args) < 2 {
		filename = "input"
	} else {
		filename = os.Args[1]
	}
	file, err := os.Open(filename)
	if err != nil {
		panic("error opening file")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	input := [][]bool{}
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]bool, len(line))
		for i, char := range line {
			if char == '@' {
				row[i] = true
			} else {
				row[i] = false
			}
		}
		input = append(input, row)
	}
	removable, _ := part1(input)
	fmt.Println(removable)
	part2(input)
}

func part1(input [][]bool) (int, [][]int8) {
	rows := len(input)
	cols := len(input[0])
	occupiedNeighbors := make([][]int8, rows)
	for i := range occupiedNeighbors {
		occupiedNeighbors[i] = make([]int8, cols)
	}
	directions := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1},
	}
	for i := range input {
		for j := range input[i] {
			if input[i][j] {
				for diff := range directions {
					di := directions[diff][0]
					dj := directions[diff][1]
					ni := i + di
					nj := j + dj
					if ni >= 0 && ni < rows && nj >= 0 && nj < cols {
						occupiedNeighbors[ni][nj]++
					}
				}
			} else {
				occupiedNeighbors[i][j] = 9
			}
		}
	}
	removable := checkArray(occupiedNeighbors, 4)
	return removable, occupiedNeighbors
}

func part2(input [][]bool) {
	sum := 0
	iters := 0
	limit := int8(4)
	removable, occupiedNeighbors := part1(input)
	for removable > 0 {
		sum += removable
		iters++
		input = updateInput(input, occupiedNeighbors, limit)
		removable, occupiedNeighbors = part1(input)
	}
	fmt.Printf("Removed %d rolls in %d iterations\n", sum, iters)
}

func updateInput(input [][]bool, occupiedNeighbors [][]int8, limit int8) [][]bool {

	for i := range input {
		for j := range input[i] {
			if input[i][j] && occupiedNeighbors[i][j] < limit {
				input[i][j] = false
			}
		}
	}
	return input
}

func checkArray(arr [][]int8, limit int8) int {
	count := 0
	for i := range arr {
		for j := range arr[i] {
			if arr[i][j] < limit {
				count++
			}
		}
	}
	return count
}
