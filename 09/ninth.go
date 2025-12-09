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
	for scanner.Scan() {
		line := scanner.Text()
		t := Tile{}
		fmt.Sscanf(line, "%d,%d", &t.x, &t.y)
		redtiles = append(redtiles, t)
	}
	edges := makeEdges(&redtiles)
	biggest := 0
	corners := make([]Tile, 2)
	//---visualizer---
	viz := NewVisualizer("day09.mp4", &edges, &redtiles)
	defer viz.Close()
	pairsChecked := 0
	renderMod := 100
	//---visualizer end---
	fmt.Println("starting to loop through candidate squares")
	for _, t1 := range redtiles {
		for _, t2 := range redtiles {

			if t1 == t2 {
				continue
			}
			dx := abs(t1.x-t2.x) + 1
			dy := abs(t1.y-t2.y) + 1
			square := dx * dy
			if square > biggest {
				pairsChecked++
				if isLegal(t1, t2, &edges) {
					corners[0] = t1
					corners[1] = t2
					biggest = square
					//---visualizer---
					for k := 0; k < 30; k++ {
						viz.AddFrame(t1, t2, corners[0], corners[1], true, biggest)
					}
					//---visualizer end---
					fmt.Println(time.Now(), "New biggest:", biggest, "Corners:", corners[0], corners[1])
				} else {
					//---visualizer---
					if pairsChecked%renderMod == 0 {
						viz.AddFrame(t1, t2, corners[0], corners[1], false, biggest)
					}
					//---visualizer end---
				}

			}
		}
	}
	fmt.Printf("Area: %d", int(biggest))
	fmt.Println("Corners:", corners[0], corners[1])
}

func isLegal(t1, t2 Tile, edges *[]Edge) bool {
	//normalize corners, so rx1, ry1 is bottom-left, rx2, ry2 is top-right
	rx1, rx2 := t1.x, t2.x
	if rx1 > rx2 {
		rx1, rx2 = rx2, rx1
	}
	ry1, ry2 := t1.y, t2.y
	if ry1 > ry2 {
		ry1, ry2 = ry2, ry1
	}
	for x := rx1; x <= rx2; x++ { // iterate over the range of x values covered by the square
		onVerticalEdge := false
		for _, edge := range *edges { // vertical edge check
			if edge.start.x == edge.end.x && edge.start.x == x {
				yMin, yMax := edge.start.y, edge.end.y
				if yMin > yMax {
					yMin, yMax = yMax, yMin
				}
				if ry1 >= yMin && ry2 <= yMax {
					onVerticalEdge = true
					break
				}
			}
		}
		if onVerticalEdge {
			continue // This column is valid, move to next x
		}
		// horizontal edge intersection check
		intersections := make([]int, 0)
		for _, edge := range *edges {
			if edge.start.y != edge.end.y { //skip the ones we already checked
				continue
			}
			ex1, ex2 := edge.start.x, edge.end.x
			if ex1 > ex2 {
				ex1, ex2 = ex2, ex1
			}
			if x >= ex1 && x < ex2 {
				intersections = append(intersections, edge.start.y)
			}
		}
		sort.Ints(intersections)
		columnValid := false
		for i := 0; i < len(intersections)-1; i += 2 {
			yBot := intersections[i]
			yTop := intersections[i+1]
			if ry1 >= yBot && ry2 <= yTop {
				columnValid = true
				break
			}
		}
		if !columnValid {
			return false //stop at invalid column
		}
	}
	return true // all columns are valid
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
func makeEdges(redtiles *[]Tile) []Edge {
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
	return edges
}
