package main

import "fmt"

func average(salary []int) float64 {

	min := salary[0]
	max := salary[0]
	total := 0
	n := len(salary)

	for i := 0; i < n; i++ {
		if salary[i] < min {
			min = salary[i]
		}
		if salary[i] > max {
			max = salary[i]
		}

		total += salary[i]
	}

	total = total - max - min
	result := float64(total) / float64(n-2)
	return result
}

func main() {
	salary := []int{1000, 2000, 4000, 3000}
	result := average(salary)
	fmt.Println(result)
}
