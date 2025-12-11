package main

import (
	"bufio"
	"fmt"
	"os"
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

		for _, conns := range connections[startNode] {
			if conns == endNode {
				continue
			}
			path := []string{}
			visited := make(map[string]bool)
			path = append(path, startNode)
			numPaths[i] += countPathsWithout(conns, connections, &path, visited, endNode, forbiddenNode)
		}
		fmt.Println("Finished subpath ", subPathDef)
	}
	for i, np := range numPaths {
		fmt.Printf("Paths from %s to %s without %s: %d\n", subPathDefs[i][0], subPathDefs[i][1], subPathDefs[i][2], np)
	}
	fmt.Printf("Total paths : %d\n", numPaths[0]*numPaths[3]*numPaths[5]+numPaths[1]*numPaths[2]*numPaths[4])
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
