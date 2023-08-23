package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Dircheck checks if the directory to write to is in the same directory as the root
// directory of the zip file. This prevents the user from writing to a directory
// outside of the zip file.
func Dircheck(root string, dest string) bool {
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return false
	}
	if len(absDest) < len(root) {
		return false
	}
	// This is a clever way that Copilot came up with to check if the destination
	// directory is in the same directory as the root directory of the zip file.
	// without actually parsing the filepaths
	for i := range root {
		if root[i] != absDest[i] {
			return false
		}
	}
	return true
}

func Dotdotcheck(dest string) bool {
	return !strings.Contains(dest, "..")
}

// Sizecheck checks if the size of the file is too large.
func Sizecheck(size int64, maxSize int64) bool {
	if size > maxSize {
		return false
	}
	return true
}

// Typecheck checks if the file is a regular file or folder
func Typecheck(file fs.FileInfo) bool {
	return file.Mode().IsDir() || file.Mode().IsRegular()
}

// Symlinksafetycheck checks if the destination of a symlink is safe
func Symlinksafetycheck(root, symlinkPath string) bool {
	dest, err := os.Readlink(symlinkPath)
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return false
	}
	if len(absDest) < len(root) {
		return false
	}
	for i := range root {
		if root[i] != absDest[i] {
			return false
		}
	}
	return true
}
