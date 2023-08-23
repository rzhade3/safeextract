package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Unzip a zip file to a destination directory
func Unzip(src, dest string, opts Options) (error, []error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err, nil
	}
	defer r.Close()

	// Create output directory if it doesn't exist
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return err, nil
	}
	if !opts.validate {
		if err := os.MkdirAll(absDest, 0755); err != nil {
			return fmt.Errorf("Error creating output directory: %v", err), nil
		}
	}

	var warnings []error
	var extractedSize int64
	// Limit the size of the extracted file
	// Iterate through the files in the archive,
	// and extract them to the destination directory
	for _, file := range r.File {
		extractedSize += int64(file.UncompressedSize64)
		if extractedSize > opts.maxSize {
			return fmt.Errorf("File size too large"), warnings
		}

		outputPath := filepath.Join(dest, file.Name)
		// Check if the destination of the file is in current directory
		if !Dircheck(absDest, outputPath) {
			warnings = append(warnings, fmt.Errorf("Destination outside of current directory: %s", outputPath))
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return err, warnings
		}
		// Check type of file
		switch file.Mode().Type() {
		case fs.ModeDir:
			if opts.validate {
				continue
			}
			// Create directory
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return fmt.Errorf("Error creating directory: %v", err), warnings
			}
		// Regular file
		case 0:
			if opts.validate {
				continue
			}
			// Create file
			fw, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("Error creating file: %v", err), warnings
			}
			defer fw.Close()
			// Copy file contents
			if _, err := io.Copy(fw, rc); err != nil {
				return fmt.Errorf("Error copying file contents: %v", err), warnings
			}
		case fs.ModeSymlink:
			if !opts.allowSymlinks {
				warnings = append(warnings, fmt.Errorf("Symlinks not allowed: %s", outputPath))
				continue
			}
			// Check if the  symlink is safe
			linkTarget, err := io.ReadAll(rc)
			if err != nil {
				return fmt.Errorf("Error reading symlink: %v", err), warnings
			}
			absLinkDest := filepath.Join(absDest, string(linkTarget))
			if !Dircheck(absDest, absLinkDest) {
				warnings = append(warnings, fmt.Errorf("Symlink destination outside of current directory: %s -> %s", file.Name, linkTarget))
				continue
			}
			if opts.validate {
				continue
			}
			// Create symlink
			if err := os.Symlink(string(linkTarget), outputPath); err != nil {
				return fmt.Errorf("Error creating symlink: %v", err), warnings
			}
		default:
			warnings = append(warnings, fmt.Errorf("File type not supported: %s", file.Name))
			continue
		}
		rc.Close()
	}
	return nil, warnings
}
