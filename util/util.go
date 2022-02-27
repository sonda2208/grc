package util

import (
	"os"
	"path/filepath"
)

var (
	commonSearchPaths = []string{
		".",
		"..",
		"../..",
		"../../..",
	}
)

func findPath(path string, searchPaths []string, filter func(os.FileInfo) bool) string {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path
		}

		return ""
	}

	for _, parent := range searchPaths {
		found, err := filepath.Abs(filepath.Join(parent, path))
		if err != nil {
			continue
		}

		fi, err := os.Stat(found)
		if err != nil {
			continue
		}

		if filter != nil {
			if filter(fi) {
				return found
			}
		} else {
			return found
		}
	}

	return ""
}

func FindFile(path string) string {
	return findPath(path, commonSearchPaths, func(fileInfo os.FileInfo) bool {
		return !fileInfo.IsDir()
	})
}

func FindDir(dir string) (string, bool) {
	found := findPath(dir, commonSearchPaths, func(fileInfo os.FileInfo) bool {
		return fileInfo.IsDir()
	})
	if found == "" {
		return "./", false
	}

	return found, true
}
