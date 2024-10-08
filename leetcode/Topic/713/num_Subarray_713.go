package main

import "fmt"

func numSubarrayProductLessThanK(nums []int, k int) int {

	// tách thành các mảng con , rồi lưu lại vào biến
	// tính tích các mảng con lại so sánh với tham số k
	// mỗi lần tích các mảng con bé hơn K thì sẽ tăng biến đếm lên

	if k <= 1 {
		return 0
	}
	count := 0
	for i := 0; i < len(nums); i++ {

	}
	return count

}

func main() {

	nums := []int{10, 5, 2, 6}
	k := 100
	result := numSubarrayProductLessThanK(nums, k)
	fmt.Println(result)

}
