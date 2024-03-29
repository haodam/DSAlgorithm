package leetcode

func CountOdds(low int, high int) int {
	result := 0
	i := low
	for i <= high {
		if i%2 == 1 {
			result++
		}
		i++
	}
	return result
}
