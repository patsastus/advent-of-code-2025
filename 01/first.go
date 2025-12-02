package main

import "bufio"
import "fmt"
import "os"


func countZeros(oldPos, newPos int, rotLeft bool)  int {
	if rotLeft  {
		oldPos--
		newPos--
	}
	start := oldPos / 100
	if oldPos < 0 && oldPos % 100 != 0 { start-- }
	end := newPos / 100
	if (newPos < 0 && newPos % 100 != 0 ) { end-- }
	if  rotLeft { return (start - end) } //how many 100s did we move left
	return (end - start) //how many hundreds did we move right
}

//part two
func main(){
	count := 0
	pos := 50
	file, err := os.Open("input.txt")
	if err != nil {panic("")}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var letter byte
		var number int
		newPos := pos
		_, err := fmt.Sscanf(line, "%c%d", &letter, &number)
		if (err != nil) {panic("")}
		if (letter == 'L') {
			newPos -= number
			count += countZeros(pos, newPos, true)
		}else if (letter == 'R') {
			newPos += number
			count += countZeros(pos, newPos, false)
		}
		pos = newPos
	}
	fmt.Print(count)
}



/*part one
func main(){
	count := 0
	pos := 50
	file, err := os.Open("input.txt")
	if err != nil {panic("")}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var letter byte
		var number int
		_, err := fmt.Sscanf(line, "%c%d", &letter, &number)
		if (err != nil) {panic("")}
		if (letter == 'L') {
			pos -= (number % 100)
		}
		if (letter == 'R') {
			pos += (number % 100)
		}
		if pos < 0 {pos += 100}
		if pos > 99 {pos = pos % 100}
		if pos == 0 {count++}
	}
	fmt.Print(count)
}
*/
