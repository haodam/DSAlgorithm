package main

import "fmt"

func canMakeArithmeticProgression(arr []int) bool {
	N := len(arr)
	min := arr[0]
	max := arr[0]
	elementExitsMap := make(map[int]bool)

	for i := 0; i < N; i++ {
		// check min in arr
		if arr[i] < min {
			min = arr[i]
		}
		// check max in arr
		if arr[i] > max {
			max = arr[i]
		}
		elementExitsMap[arr[i]] = true
	}

	totalDiff := max - min
	if totalDiff%(N-1) != 0 {
		return false
	}

	diff := totalDiff / (N - 1)
	for min < max {
		min += diff
		_, isExits := elementExitsMap[min]
		if !isExits {
			return false
		}
	}

	return true
}

func main() {

	arr := []int{1, 3, 5, 7, 9, 11, 13, 15}
	result := canMakeArithmeticProgression(arr)
	fmt.Println(result)
}
