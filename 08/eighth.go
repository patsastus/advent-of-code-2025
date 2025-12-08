package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . <filename> <count>")
		os.Exit(1)
	}
	filename := os.Args[1]
	count, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("")
	}
	file, err := os.Open(filename)
	if err != nil {
		panic("")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	partOne(scanner, count)
	partTwo(scanner)
	fmt.Println("Elapsed time:", time.Since(start))
}

type Junction struct {
	x, y, z, id          int
	nearest, circuitHead *Junction
	distance             float64
}

func calcNearest(circuits *map[*Junction][]*Junction) {
	for junc1 := range *circuits {
		distance := math.MaxFloat64
		
		for junc2 := range *circuits {
			if junc1.id != junc2.id {
				dx := junc1.x - junc2.x
				dy := junc1.y - junc2.y
				dz := junc1.z - junc2.z
				d2 := float64(dx*dx + dy*dy + dz*dz)
				if d2 < distance {
					distance = d2
					junc1.nearest = junc2
					junc1.distance = d2
				}
			}
		}
	}
}

func onCircuit(junc1, junc2 *Junction, circuits *map[*Junction][]*Junction) bool {
	if junc1.circuitHead == junc2.circuitHead {
		return true
	}
	found := false
	for _, junc := range (*circuits)[junc1.circuitHead] {
		if junc.id == junc2.id {
			found = true
			break
		}
	}
	for _, junc := range (*circuits)[junc2.circuitHead] {
		if junc.id == junc1.id {
			found = true
			break
		}
	}
	return found
}

func joinCircuits(junc1, junc2 *Junction, circuits *map[*Junction][]*Junction) {
	head1 := junc1.circuitHead
	head2 := junc2.circuitHead
	(*circuits)[head1] = append((*circuits)[head1], head2)
	(*circuits)[head1] = append((*circuits)[head1], (*circuits)[head2]...)
	junc2.circuitHead = head1
	for _, junc := range (*circuits)[head2] {
		junc.circuitHead = head1
	}
	delete(*circuits, head2)
}

func partOne(scanner *bufio.Scanner, count int) int {
	circuits := make(map[*Junction][]*Junction)
	id := 0
	for scanner.Scan() {
		line := scanner.Text()
		newJunction := &Junction{}
		fmt.Scanf(line, "%d,%d,%d", &newJunction.x, &newJunction.y, &newJunction.z)
		newJunction.id = id
		newJunction.circuitHead = newJunction
		id++
		circuits[newJunction] = []*Junction{}
	}
	calcNearest(&circuits)
	for i := 0; i < count; i++ {
		// Find the shortest distance junction
		minDistance := math.MaxFloat64
		var candidate, nearest *Junction
		for junc := range circuits {
			if junc.distance < minDistance && !onCircuit(junc, junc.nearest, &circuits) {
				minDistance = junc.distance
				candidate = junc
				nearest = junc.nearest
			}
		}
		// Join the circuits
		fmt.Printf("joining %v and %v\n", candidate, nearest)
		joinCircuits(candidate, nearest, &circuits)
	}

	return 0
}

func partTwo(scanner *bufio.Scanner) int {
	return 0
}
