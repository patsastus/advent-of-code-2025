package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"time"
)

type Tile struct {
	x, y int
}

type Edge struct {
	start, end Tile
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
	//partOne(scanner)
	partTwo(scanner)
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

func partTwo(scanner *bufio.Scanner) {
	redtiles := make([]Tile, 0)
	greentiles := make(map[Tile]bool)
	for scanner.Scan() {
		line := scanner.Text()
		t := Tile{}
		fmt.Sscanf(line, "%d,%d", &t.x, &t.y)
		redtiles = append(redtiles, t)
		if len(redtiles) > 1 {
			prev := redtiles[len(redtiles)-2]
			if t.x == prev.x {
				for y := min(t.y, prev.y) + 1; y < max(t.y, prev.y); y++ {
					greentiles[Tile{t.x, y}] = true
				}
			} else if t.y == prev.y {
				for x := min(t.x, prev.x) + 1; x < max(t.x, prev.x); x++ {
					greentiles[Tile{x, t.y}] = true
				}
			}
		}
	}
	fillGreenTiles(&redtiles, &greentiles)
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

func fillGreenTiles(redtiles *[]Tile, greentiles *map[Tile]bool) {
	edges := make([]Edge, 0)
	previousCorner := (*redtiles)[len(*redtiles)-1]
	minX, maxX := previousCorner.x, previousCorner.x
	for _, corner := range *redtiles {
		edges = append(edges, Edge{previousCorner, corner})
		if corner.x < minX {
			minX = corner.x
		}
		if corner.x > maxX {
			maxX = corner.x
		}
		previousCorner = corner
	}
	sort.Slice(edges, func(i, j int) bool { //sort by x, then by y
		if edges[i].start.x != edges[j].start.x {
			return edges[i].start.x < edges[j].start.x
		}
		return edges[i].start.y < edges[j].start.y
	})
	for x := minX; x <= maxX; x++ { // iterate over the range of x values covered by the edges
		intersections := make([]int, 0) //at this x, what y values do we intersect with edges
		for _, edge := range edges {
			x1, x2 := edge.start.x, edge.end.x
			if x1 == x2 {
				continue
			} // vertical edge, skip it
			if x1 > x2 {
				x1, x2 = x2, x1
			} // swap if necessary
			if x >= x1 && x <= x2 {
				intersections = append(intersections, edge.start.y)
			}
		}
		sort.Ints(intersections)
		for i := 0; i < len(intersections)-1; i += 2 { //every other range is filled
			yStart := intersections[i]
			yEnd := intersections[i+1]
			for y := yStart; y <= yEnd; y++ {
				(*greentiles)[Tile{x, y}] = true
			}
		}
	}

}
