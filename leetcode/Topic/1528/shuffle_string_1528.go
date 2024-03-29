package main

import "fmt"

func RestoreString(s string, indices []int) string {

	result := make([]byte, len(indices))
	for index, value := range indices {
		result[value] = s[index]
	}

	return string(result)

}

func main() {

	s := "codeleet"
	indices := []int{4, 5, 6, 7, 0, 2, 1, 3}

	result := RestoreString(s, indices)

	fmt.Println(result)
}
