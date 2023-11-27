package main

import (
	"fmt"
	leetcode "github.com/haodam/DSAlgorithm/leetcode/1512"
)

func main() {
	//var temp = leetcode.NumWaterBottles(9, 3)
	//fmt.Println(temp)

	array := [...]int{1, 2, 3, 1, 1, 3}
	var temp = leetcode.NumIdenticalPairs(array[:])
	fmt.Println(temp)

}
