package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
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
	partOne(scanner)
	fmt.Println("Elapsed time:", time.Since(start))
}

func partOne(scanner *bufio.Scanner) {
	connections := make(map[string][]string)
	for scanner.Scan() {
		line := scanner.Text()
		firstSplit := strings.Split(line, ":")
		connections[firstSplit[0]] = strings.Fields(firstSplit[1])
	}
	num := 0
	visited := make(map[string]bool)
	visited["you"] = true
	for _, conns := range connections["you"] {
		if conns == "out" {
			num++
			//		fmt.Printf("found path via %v\n", visited)
			continue
		}
		num += countPaths(conns, connections, visited)
	}
	fmt.Println("Total paths from 'you' to 'out':", num)
	generateGraphviz(connections)
	subPathDefs := [][]string{
		{"svr", "fft", "dac"},
		{"svr", "dac", "fft"},
		{"dac", "fft", ""},
		{"fft", "dac", ""},
		{"fft", "out", "dac"},
		{"dac", "out", "fft"},
	}
	numPaths := make([]int, len(subPathDefs))
	for i, subPathDef := range subPathDefs {
		startNode := subPathDef[0]
		endNode := subPathDef[1]
		forbiddenNode := subPathDef[2]

		for _, conns := range connections[startNode] {
			if conns == endNode {
				continue
			}
			path := []string{}
			memo := make(map[string]int)
			path = append(path, startNode)
			numPaths[i] += countPathsMemo(conns, connections, &path, memo, endNode, forbiddenNode)
		}
	}
	fmt.Printf("Total paths : %d\n", numPaths[0]*numPaths[3]*numPaths[5]+numPaths[1]*numPaths[2]*numPaths[4])
}

func countPathsMemo(node string, connections map[string][]string, pathSoFar *[]string, memo map[string]int, target, forbidden string) int {
	if val, ok := memo[node]; ok {
		return val
	}
	*pathSoFar = append(*pathSoFar, node)
	defer func() {
		*pathSoFar = (*pathSoFar)[:len(*pathSoFar)-1]
	}()
	if node == target {
		return 1
	}
	if node == forbidden {
		return 0
	}
	count := 0
	for _, conn := range connections[node] {
		count += countPathsMemo(conn, connections, pathSoFar, memo, target, forbidden)
	}
	memo[node] = count
	return count
}

func countPathsWithout(node string, connections map[string][]string, pathSoFar *[]string, visited map[string]bool, target, forbidden string) int {
	count := 0
	*pathSoFar = append(*pathSoFar, node)
	visited[node] = true
	defer func() { //undo these changes on return
		delete(visited, node)
		*pathSoFar = (*pathSoFar)[:len(*pathSoFar)-1]
	}()
	for _, conn := range connections[node] {
		if conn == forbidden || visited[conn] {
			continue
		}
		if conn == target {
			count++
			continue
		}
		count += countPathsWithout(conn, connections, pathSoFar, visited, target, forbidden)
	}
	return count
}

func countPaths(node string, connections map[string][]string, visited map[string]bool) int {
	count := 0
	visited[node] = true
	for _, conn := range connections[node] {
		if conn == "out" {
			count++
			//		fmt.Printf("found path via %v\n", visited)
			continue
		}
		count += countPaths(conn, connections, visited)
	}
	return count
}

func generateGraphviz(connections map[string][]string) {
	f, err := os.Create("graph.dot")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString("strict digraph {\n")

	w.WriteString("  layout=fdp\n")
	w.WriteString("  overlap=true\n")
	w.WriteString("  splines=false\n")
	w.WriteString("  node [shape=circle, style=filled, color=lightblue];\n")

	w.WriteString("  svr [pos=\"0,50!\", fillcolor=green, shape=doublecircle];\n")

	w.WriteString("  fft [pos=\"50,100!\", fillcolor=yellow, shape=doublecircle];\n")

	w.WriteString("  dac [pos=\"50,0!\", fillcolor=orange, shape=doublecircle];\n")

	w.WriteString("  out [pos=\"100,50!\", fillcolor=red, shape=doublecircle];\n")

	seen := make(map[string]bool)
	for src, targets := range connections {
		for _, dst := range targets {
			// Edge deduplication logic
			a, b := src, dst
			if a > b {
				a, b = b, a
			}
			edgeKey := a + "-" + b

			if !seen[edgeKey] {
				w.WriteString(fmt.Sprintf("  %s -> %s;\n", src, dst))
				seen[edgeKey] = true
			}
		}
	}
	w.WriteString("}\n")
	w.Flush()
	fmt.Println("Generated graph.dot with pinned nodes.")
}
