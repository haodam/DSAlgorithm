package scanner

import "fmt"

func FormatSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size >= GB:
		return fmt.Sprintf(
			"%.2f GB",
			float64(size)/GB,
		)

	case size >= MB:
		return fmt.Sprintf(
			"%.2f MB",
			float64(size)/MB,
		)

	case size >= KB:
		return fmt.Sprintf(
			"%.2f KB",
			float64(size)/KB,
		)

	default:
		return fmt.Sprintf(
			"%d B",
			size,
		)
	}
}
