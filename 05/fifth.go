package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	in := [][]int{}
	mode := false
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			mode = true
			continue
		}
		if !mode {
			inRange := make([]int, 2)
			fmt.Sscanf(line, "%d-%d", &inRange[0], &inRange[1])
			in = append(in, inRange)
		} else {
			num := 0
			fmt.Sscanf(line, "%d", &num)
			found := false
			for _, r := range in {
				if num >= r[0] && num <= r[1] {
					found = true
					break
				}
			}
			if found {
				count++
			}
		}
	}
	fmt.Println(count)
	sort.Slice(in, func(i, j int) bool {
		return in[i][0] < in[j][0]
	})
	fmt.Println(in)
	collapsed := collapseRanges(in)
	fmt.Print(collapsed)
	fmt.Println(sumRanges(collapsed))
}

func collapseRanges(ranges [][]int) [][]int {
	collapsed := [][]int{}
	current := ranges[0]
	for i := 1; i < len(ranges); i++ {
		if ranges[i][0] <= current[1] {
			if ranges[i][1] > current[1] {
				current[1] = ranges[i][1]
			}
		} else {
			collapsed = append(collapsed, current)
			current = ranges[i]
		}
	}
	collapsed = append(collapsed, current)
	return collapsed
}

func sumRanges(ranges [][]int) int {
	sum := 0
	for _, r := range ranges {
		sum += r[1] - r[0] + 1
	}
	return sum
}
