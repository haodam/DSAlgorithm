package main

import (
	"fmt"
	"math"
)

func containsNearbyDuplicate(nums []int, k int) bool {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] && math.Abs(float64(i-j)) <= float64(k) {
				return true
			}
		}
	}
	return false
}

func containsNearbyDuplicateV2(nums []int, k int) bool {
	seen := make(map[int]int)
	for i, num := range nums {
		if idx, ok := seen[num]; ok && i-idx <= k {
			return true
		}
		seen[num] = i
	}
	return false
}

func main() {
	nums := []int{1, 2, 3, 1}
	k := 3
	fmt.Println(containsNearbyDuplicate(nums, k))
	fmt.Println(containsNearbyDuplicateV2(nums, k))
}
