package main

import (
	"bufio"
	"fmt"
	"os"
)

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
	//partOne(scanner)
	partTwo(scanner)
}

func partOne(scanner *bufio.Scanner) {
	count := 0
	previousLine := []rune(scanner.Text())
	//fmt.Println(string(previousLine))
	for scanner.Scan() {
		line := []rune(scanner.Text())
		//		fmt.Println(string(line))
		for i := 0; i < len(previousLine); i++ {
			if previousLine[i] == 'S' || previousLine[i] == '|' {
				if line[i] == '.' {
					line[i] = '|'
				}
				if line[i] == '^' {
					count++
					if i > 0 && line[i-1] == '.' {
						line[i-1] = '|'
					}
					if i < len(line)-1 && line[i+1] == '.' {
						line[i+1] = '|'
					}
				}
			}
		}
		previousLine = line
	}
	fmt.Println(count)
	fmt.Println(string(previousLine))
}

func lineToPathCount(line []rune) []int {

	path := make([]int, len(line))
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case 'S':
			path[i] = 1
		case '^':
			path[i] = -1
		default:
			path[i] = 0
		}
	}
	//	fmt.Println(path)
	return path
}
func partTwo(scanner *bufio.Scanner) {
	count := 1
	scanner.Scan()
	previousLine := []rune(scanner.Text())
	previousPaths := lineToPathCount(previousLine)
	for scanner.Scan() {
		line := []rune(scanner.Text())
		path := lineToPathCount(line)
		for i := 0; i < len(path); i++ {
			if path[i] >= 0 { //splitters are the only negative values, others get what's above added
				path[i] += max(0, previousPaths[i])
			} else if path[i] < 0 { //splitters add what's coming from above to either side, if free and existing
				if i > 0 && path[i-1] >= 0 {
					path[i-1] += max(0, previousPaths[i])
				}
				if i < len(path)-1 && path[i+1] >= 0 {
					path[i+1] += max(0, previousPaths[i])
				}
			}
		}
		previousLine = line
		previousPaths = path
		//		fmt.Printf("%v\n", path)
		count++
	}
	sum := 0
	for _, v := range previousPaths {
		sum += v
	}
	fmt.Println(sum)
}
