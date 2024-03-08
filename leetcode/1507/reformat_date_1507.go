package main

import "fmt"

func reformatDate(date string) string {
	result := ""
	monthMap := map[string]string{
		"Jan": "01",
		"Feb": "02",
		"Mar": "03",
		"Apr": "04",
		"May": "05",
		"Jun": "06",
		"Jul": "07",
		"Aug": "08",
		"Sep": "09",
		"Oct": "10",
		"Nov": "11",
		"Dec": "12",
	}

	// chieu dai toi da cua chuoi date la 13 ky tu
	if len(date) == 12 {
		date = "0" + date
	}

	year := date[9:]
	month := monthMap[date[5:8]]
	day := date[:2]
	result = year + "-" + month + "-" + day

	return result
}

func main() {
	date := "20th Oct 2052"
	result := reformatDate(date)
	fmt.Println(result)
}
