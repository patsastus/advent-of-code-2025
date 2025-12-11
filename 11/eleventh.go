package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
	//	"math"
	//	"sort"
)

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
	subPathDefs := [][]string{
		{"srv", "fft", "dac"},
		{"srv", "dac", "fft"},
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
		path := []string{}
		path = append(path, startNode)
		for _, conns := range connections[startNode] {
			if conns == endNode {
				continue
			}
			copied := make([]string, len(path))
			copy(copied, path)
			numPaths[i] += countPathsWithout(conns, connections, copied, endNode, forbiddenNode)
		}
	}
	for i, np := range numPaths {
		fmt.Printf("Paths from %s to %s without %s: %d\n", subPathDefs[i][0], subPathDefs[i][1], subPathDefs[i][2], np)
	}
	fmt.Printf("Total paths : %d\n", numPaths[0]*numPaths[3]*numPaths[5]+numPaths[1]*numPaths[2]*numPaths[4])
}

func countPathsWith(node string, connections map[string][]string, visited []string) int {
	count := 0
	visited = append(visited, node)
	for _, conn := range connections[node] {
		if conn == "out" {
			if slices.Contains(visited, "dac") && slices.Contains(visited, "fft") {
				count++
			}
			//	fmt.Printf("found path via %v\n", visited)
			continue
		}
		copied := make([]string, len(visited))
		copy(copied, visited)
		count += countPathsWith(conn, connections, copied)
	}
	return count
}

func countPathsWithout(node string, connections map[string][]string, visited []string, target, forbidden string) int {
	count := 0
	visited = append(visited, node)
	for _, conn := range connections[node] {
		if conn == forbidden {
			continue
		}
		if conn == target {
			count++
			//	fmt.Printf("found path via %v\n", visited)
			continue
		}
		copied := make([]string, len(visited))
		copy(copied, visited)
		count += countPathsWith(conn, connections, copied)
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
