package main

import (
	"fmt"
	"strconv"
)

func isPathCrossing(path string) bool {
	x := 0
	y := 0
	visitedNode := make(map[string]bool)
	visitedNode["00"] = true

	for _, p := range path {
		switch p {
		case 'N':
			y++
		case 'S':
			y--
		case 'E':
			x++
		case 'W':
			x--
		}

		key := strconv.Itoa(x) + strconv.Itoa(y)
		_, visited := visitedNode[key]
		if visited {
			return true
		} else {
			visitedNode[key] = true
		}
	}
	return false
}

func main() {
	path := "NES"
	result := isPathCrossing(path)
	fmt.Println(result)
}
