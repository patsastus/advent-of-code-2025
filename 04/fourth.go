package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <filename>")
        os.Exit(1)
    }
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {panic("")}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	input := [][]bool{}
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]bool, len(line))
		for i, char := range line {
			if char == '@' { row[i] = true } else { row[i] = false }
		}
		input = append(input, row)
	}
	part1(input)
	part2(input)
}	

func part1(input [][]bool) [][]int8 {
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
	fmt.Print(removable)
	return occupiedNeighbors
}

func part2(input [][]bool) {
	sum := 0
	iters := 0
	limit := int8(4)
	occupiedNeighbors := part1(input)
	removable := checkArray(occupiedNeighbors, limit)
	for removable > 0 {
		sum += removable
		iters++
		input = updateInput(input, occupiedNeighbors, limit)
		occupiedNeighbors = part1(input)
		removable = checkArray(occupiedNeighbors, limit)
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
			//if arr[i][j] < 9 {fmt.Print(arr[i][j])} else {fmt.Print(".")}
			if(arr[i][j] < limit) {
				count++
			}
		}
		fmt.Println()
	}
	return count
}