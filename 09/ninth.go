package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

type Tile struct {
	x, y int
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
	//partTwo(scanner)
	fmt.Println("Elapsed time:", time.Since(start))
}

func partOne(scanner *bufio.Scanner) {
	redtiles := make(map[Tile]bool)
	for scanner.Scan() {
		line := scanner.Text()
		t := Tile{}
		fmt.Sscanf(line, "%d,%d", &t.x, &t.y)
		redtiles[t] = true
	}
	biggest := 0.0
	corners := make([]Tile, 2)
	for t1 := range redtiles {
		for t2 := range redtiles {
			if t1 == t2 {
				continue
			}
			dx := math.Abs(float64(t1.x-t2.x)) + 1
			dy := math.Abs(float64(t1.y-t2.y)) + 1
			square := dx * dy
			if square > biggest {
				corners[0] = t1
				corners[1] = t2
				biggest = square
			}
		}
	}
	fmt.Printf("Area: %d", int(biggest))
	fmt.Println("Corners:", corners[0], corners[1])
}
