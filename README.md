# Advent of Code 2025

[**Advent of Code**](https://adventofcode.com/) is a yearly coding challenge, where you get a small puzzle to solve algorithmically. Typically, the problems are published daily, starting on December 1st, and the puzzles generally involve some parsing of input text, choosing the right data structure to make your algorithm feasible, and some algorithm design in how you go about solving the puzzle. 

## 2025 version
This was my first go at taking part in **AoC**, and there was a slight departure from previous years: instead of 25 days of individual puzzles, there would be 12 days of 2-part puzzles, where the second part uses the same input as the first, but has a different puzzle to solve. I chose to use this as a chance to sharpen my **Go** skills, which I'm still a little new at. I managed to finish each puzzle the day it released, although Day 10 took an embarassing amount of my "workday" that day.

## Prerequisites
All other days use only Go standard library functions, but for Day 10, I needed a Liner Programming solver, which I couldn't find in Go, and I really didn't think I could do it myself within a day. So I found a Go wrapper of the C library **GLPK**, which needs to be installed:
```
sudo apt-get install libglpk-dev
```

## Testing
Each days' solution is in its own folder. The **AoC** website asks you not to include your input file in any public repos, so to test these out, you will need to go to [**Advent of Code**](https://adventofcode.com/) and generate an input file for that day, and save it in the days folder as `input`. Running
```
go run .
```
will then output the answers to the two puzzles, which you can check on the website. 

## Extra visuals
For days 8 and 9 I made video visualizers to try to better wrap my head around what was happening in the algorithm I wrote, and for day 11 I made a graph visualization of the input data. I've added the visuals created here (Any ones you generate will differ slightly in detail because of different input)

### Day 8: connecting 3D dots to the closest neighbor

### Day 9: Find the largest square that doesn't cross a line

### Day 11: DAG
