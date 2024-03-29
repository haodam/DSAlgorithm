package leetcode

func NumIdenticalPairs(nums []int) int {
	var count = 0
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums)-1; j++ {
			if nums[i] == nums[j] {
				count++
			}
		}
	}
	return count
}

func NumIdenticalPairsV2(nums []int) int {
	var count = 0
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums)-1; j++ {
			if nums[i] == nums[j] {
				count++
			}
		}
	}
	return count
}
