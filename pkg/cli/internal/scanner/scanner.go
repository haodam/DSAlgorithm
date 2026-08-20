package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type FileInfo struct {
	Path string
	Size int64
}

type ScanResult struct {
	Files     []FileInfo
	TotalSize int64
}

func Scan(path string) (*ScanResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot access path %q: %w",
			path,
			err,
		)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf(
			"%q is not a directory",
			path,
		)
	}

	result := &ScanResult{
		Files: make([]FileInfo, 0),
	}

	err = filepath.Walk(path, func(
		currentPath string,
		info os.FileInfo,
		err error,
	) error {
		if err != nil {
			return fmt.Errorf(
				"cannot access %q: %w",
				currentPath,
				err,
			)
		}

		// Directory không cần đưa vào result.
		if info.IsDir() {
			return nil
		}

		file := FileInfo{
			Path: currentPath,
			Size: info.Size(),
		}

		result.Files = append(
			result.Files,
			file,
		)

		result.TotalSize += info.Size()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// File lớn nhất đứng trước.
	sort.Slice(result.Files, func(i, j int) bool {
		return result.Files[i].Size > result.Files[j].Size
	})

	return result, nil
}
