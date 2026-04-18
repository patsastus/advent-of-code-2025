package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	//	"math"
	//	"slices"
	//	"sort"

	"github.com/lukpank/go-glpk/glpk"
)

type Machine struct {
	lights  []bool
	buttons [][]int
	joltage []int
	ID      int
}

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
	scanner := bufio.NewScanner(file)
	partOne(scanner)
	file.Close()
	file, err = os.Open(filename)
	if err != nil {
		panic("error opening file")
	}
	scanner = bufio.NewScanner(file)
	partTwo(scanner)
	fmt.Println("Elapsed time:", time.Since(start))
}

func partOne(scanner *bufio.Scanner) {
	machines := make([]Machine, 0)
	for scanner.Scan() {
		m := Machine{}
		line := scanner.Text()
		splits := strings.Split(line, " ")
		m.lights = make([]bool, len(splits[0])-2)
		m.buttons = make([][]int, len(splits)-2)
		for i, r := range []rune(splits[0])[1 : len(splits[0])-1] {
			if r == '#' {
				m.lights[i] = true
			}
		}
		for i, btnStr := range splits[1 : len(splits)-1] {
			btnSplits := strings.Split(btnStr[1:len(btnStr)-1], ",")
			btns := make([]int, len(btnSplits))
			for j, bs := range btnSplits {
				fmt.Sscanf(bs, "%d", &btns[j])
			}
			m.buttons[i] = btns
		}
		machines = append(machines, m)
		//	fmt.Printf("Parsed machine: %+v\n", m)
	}
	sum := 0
	for _, m := range machines {
		done := false
		for n := 1; n <= len(m.buttons) && !done; n++ {
			combos := makeCombinations(n, m)
			for _, c := range combos {
				if isCorrect(m, c) {
					//	fmt.Printf("Combo %v of length %d is correct for machine %v\n", c, n, m)
					done = true
					sum += n
					break
				}
			}
		}
	}
	fmt.Printf("total combos: %d\n", sum)
}

func partTwo(scanner *bufio.Scanner) {
	machines := make([]Machine, 0)
	id := 0
	for scanner.Scan() {
		m := Machine{}
		line := scanner.Text()
		splits := strings.Split(line, " ")
		m.lights = make([]bool, len(splits[0])-2)
		m.buttons = make([][]int, len(splits)-2)
		for i, r := range []rune(splits[0])[1 : len(splits[0])-1] {
			if r == '#' {
				m.lights[i] = true
			}
		}
		for i, btnStr := range splits[1 : len(splits)-1] {
			btnSplits := strings.Split(btnStr[1:len(btnStr)-1], ",")
			btns := make([]int, len(btnSplits))
			for j, bs := range btnSplits {
				fmt.Sscanf(bs, "%d", &btns[j])
			}
			m.buttons[i] = btns
		}
		jStr := splits[len(splits)-1]
		jolts := strings.Split(jStr[1:len(jStr)-1], ",")
		m.joltage = make([]int, len(jolts))
		for i, str := range jolts {
			fmt.Sscanf(str, "%d", &m.joltage[i])
		}
		m.ID = id
		id++
		machines = append(machines, m)
	}
	sum := 0
	for _, m := range machines {
		result := solveWithGLPK(m)
		sum += sumSlice(result)
	}
	fmt.Printf("total combos: %d\n", sum)
}

func sumSlice(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func isCorrect(m Machine, c []int) bool {
	temp := make([]bool, len(m.lights))
	for _, btn := range c { //push all the buttons the combo tells you to
		for _, i := range m.buttons[btn] { //change status of all indexes on the button
			temp[i] = !temp[i]
		}
	}
	for i, _ := range temp {
		if temp[i] != m.lights[i] {
			return false
		}
	}
	return true
}

func makeCombinations(n int, m Machine) [][]int {
	var result [][]int
	indexes := make([]int, n)
	numBtns := len(m.buttons)
	i := 0
	for i >= 0 {
		combo := make([]int, n)
		for j := 0; j < n; j++ {
			combo[j] = indexes[j]
		}
		result = append(result, combo)
		i = n - 1
		for i >= 0 && indexes[i] == numBtns-1 {
			i--
		}
		if i < 0 {
			break
		}
		indexes[i]++
		for j := i + 1; j < n; j++ {
			indexes[j] = indexes[i]
		}
	}
	return result
}

func solveWithGLPK(m Machine) []int {
	lp := glpk.New()
	defer lp.Delete()

	lp.SetProbName(fmt.Sprintf("M%d", m.ID))
	lp.SetObjDir(glpk.MIN) // Minimize

	numButtons := len(m.buttons)
	numTargets := len(m.joltage)

	lp.AddCols(numButtons)
	for i := 0; i < numButtons; i++ {
		colIdx := i + 1 // GLPK is 1-based

		lp.SetColName(colIdx, fmt.Sprintf("btn_%d", i))

		lp.SetColKind(colIdx, glpk.IV)
		lp.SetColBnds(colIdx, glpk.LO, 0.0, 0.0)

		lp.SetObjCoef(colIdx, 1.0)
	}

	lp.AddRows(numTargets)

	nnz := 0
	for _, btn := range m.buttons {
		nnz += len(btn)
	}

	indRow := make([]int32, nnz+1)
	indCol := make([]int32, nnz+1)
	val := make([]float64, nnz+1)

	currentIdx := 1
	for btnIdx, affectedTargets := range m.buttons {
		for _, targetIdx := range affectedTargets {
			indRow[currentIdx] = int32(targetIdx + 1) // Row (Target)
			indCol[currentIdx] = int32(btnIdx + 1)    // Col (Button)
			val[currentIdx] = 1.0
			currentIdx++
		}
	}

	for i := 0; i < numTargets; i++ {
		rowIdx := i + 1
		targetVal := float64(m.joltage[i])

		lp.SetRowName(rowIdx, fmt.Sprintf("target_%d", i))
		lp.SetRowBnds(rowIdx, glpk.FX, targetVal, targetVal)
	}

	lp.LoadMatrix(indRow, indCol, val)

	iocp := glpk.NewIocp()
	iocp.SetPresolve(true)       // Helper to simplify problem
	iocp.SetMsgLev(glpk.MSG_OFF) // Silence output

	err := lp.Intopt(iocp)
	if err != nil {
		return nil
	}

	stat := lp.MipStatus()

	if stat != glpk.OPT && stat != glpk.FEAS {
		return nil // Impossible
	}

	counts := make([]int, numButtons)
	for i := 0; i < numButtons; i++ {
		counts[i] = int(lp.MipColVal(i + 1))
	}
	return counts
}
