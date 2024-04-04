package main

import "fmt"

func findLHS(nums []int) int {
	freq := make(map[int]int)
	for _, num := range nums {
		freq[num]++
	}

	maxLen := 0
	for num, count := range freq {
		if freq[num+1] > 0 {
			currentLen := count + freq[num+1]
			if currentLen > maxLen {
				maxLen = currentLen
			}
		}
	}
	return maxLen
}

func main() {
	nums := []int{1, 3, 2, 2, 5, 2, 3, 7}
	fmt.Println("Input nums:", nums)
	fmt.Println("Output:", findLHS(nums))
}
