package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	if (count % 2) == 0 {
		fmt.Println("boo")
	}
	//one := partOne(scanner, count)
	//fmt.Println(one)
	two := partTwo(scanner)
	fmt.Println(two)
	fmt.Println("Elapsed time:", time.Since(start))
}

type Junction struct {
	x, y, z, id int
	circuitHead *Junction
}

type Connection struct {
	from, to *Junction
	distance float64
}

func dist(junc1, junc2 *Junction) float64 {
	dx := junc1.x - junc2.x
	dy := junc1.y - junc2.y
	dz := junc1.z - junc2.z
	return float64(dx*dx + dy*dy + dz*dz)
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
	head2.circuitHead = head1
	for _, junc := range (*circuits)[head2] {
		junc.circuitHead = head1
	}
	delete(*circuits, head2)
}

func partOne(scanner *bufio.Scanner, count int) int {
	circuits := make(map[*Junction][]*Junction)
	var allJunctions []*Junction
	id := 0
	for scanner.Scan() {
		line := scanner.Text()
		newJunction := &Junction{}
		fmt.Sscanf(line, "%d,%d,%d", &newJunction.x, &newJunction.y, &newJunction.z)
		newJunction.id = id
		newJunction.circuitHead = newJunction
		id++
		circuits[newJunction] = []*Junction{}
		allJunctions = append(allJunctions, newJunction)
	}
	var connections []Connection
	for i := 0; i < len(allJunctions); i++ {
		for j := i + 1; j < len(allJunctions); j++ {
			p1 := allJunctions[i]
			p2 := allJunctions[j]
			d := dist(p1, p2)
			connections = append(connections, Connection{from: p1, to: p2, distance: d})
		}
	}
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].distance < connections[j].distance
	})
	joins := 0
	for _, c := range connections {
		if joins >= count {
			break
		}
		if c.from.circuitHead != c.to.circuitHead {
			joinCircuits(c.from, c.to, &circuits)
			fmt.Printf("Joining %v and %v at distance %v\n", c.from, c.to, c.distance)
		}
		joins++
	}
	fmt.Printf("Final %d circuits:\n", len(circuits))
	fmt.Print(circuits)
	keys := make([]*Junction, 0, len(circuits))
	for k := range circuits {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		len1 := len(circuits[keys[i]])
		len2 := len(circuits[keys[j]])
		return len1 > len2
	})
	val1, val2, val3 := 0, 0, 0
	if len(keys) > 0 {
		val1 = len(circuits[keys[0]]) + 1
	}
	if len(keys) > 1 {
		val2 = len(circuits[keys[1]]) + 1
	}
	if len(keys) > 2 {
		val3 = len(circuits[keys[2]]) + 1
	}

	fmt.Printf("Top 3 sizes: %d, %d, %d\n", val1, val2, val3)
	return val1 * val2 * val3
}

func partTwo(scanner *bufio.Scanner) int {
	circuits := make(map[*Junction][]*Junction)
	var allJunctions []*Junction
	id := 0
	for scanner.Scan() {
		line := scanner.Text()
		newJunction := &Junction{}
		fmt.Sscanf(line, "%d,%d,%d", &newJunction.x, &newJunction.y, &newJunction.z)
		newJunction.id = id
		newJunction.circuitHead = newJunction
		id++
		circuits[newJunction] = []*Junction{}
		allJunctions = append(allJunctions, newJunction)
	}
	var connections []Connection
	for i := 0; i < len(allJunctions); i++ {
		for j := i + 1; j < len(allJunctions); j++ {
			p1 := allJunctions[i]
			p2 := allJunctions[j]
			d := dist(p1, p2)
			connections = append(connections, Connection{from: p1, to: p2, distance: d})
		}
	}
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].distance < connections[j].distance
	})
	var joiner, joinee *Junction
	for _, c := range connections {
		if len(circuits) == 1 {
			break
		}
		if c.from.circuitHead != c.to.circuitHead {
			joinCircuits(c.from, c.to, &circuits)
			joiner, joinee = c.from, c.to
			fmt.Printf("Joining %v and %v at distance %v\n", c.from, c.to, c.distance)
		}
	}
	fmt.Printf("Final join %v to %v:\n", joiner, joinee)
	return joiner.x * joinee.x
}
