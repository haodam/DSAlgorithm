package main

import "fmt"

func risky() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Đã recover panic:", r)
		}
	}()
	fmt.Println("Gây panic")
	panic("Lỗi nghiêm trọng")

	fmt.Println("Không bao giờ in dòng này")
}

func main() {
	risky()
	fmt.Println("Chương trình vẫn tiếp tục chạy sau panic")
}
