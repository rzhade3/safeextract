package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Untar a tar.gz file to a destination directory
func Untar(src, dest string, opts Options) (error, []error) {
	// Open the tarball for reading.
	fr, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Error opening tarball: %v", err), nil
	}
	defer fr.Close()
	gzr, err := gzip.NewReader(fr)
	if err != nil {
		return fmt.Errorf("Error creating gzip reader: %v", err), nil
	}
	defer gzr.Close()

	absDest, err := filepath.Abs(dest)
	if err != nil {
		return err, nil
	}
	if !opts.validate {
		// Create output directory if it doesn't exist
		if err := os.MkdirAll(absDest, 0755); err != nil {
			return fmt.Errorf("Error creating output directory: %v", err), nil
		}
	}

	tr := tar.NewReader(gzr)
	var extractedSize int64

	var warnings []error
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading tar header: %v", err), warnings
		}
		extractedSize += header.Size
		// Check if the size of the extracted file is too large
		if extractedSize > opts.maxSize {
			return fmt.Errorf("File size too large"), warnings
		}
		// Check if the destination of the file is in current directory
		outputPath := filepath.Join(absDest, header.Name)
		if !Dircheck(absDest, outputPath) {
			warnings = append(warnings, fmt.Errorf("Destination outside of current directory: %s", outputPath))
			continue
		}
		// Check if the item is a file or directory
		switch header.Typeflag {
		case tar.TypeDir:
			if opts.validate {
				continue
			}
			// Create directory
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return fmt.Errorf("Error creating directory: %v", err), warnings
			}
		case tar.TypeReg:
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
			if _, err := io.Copy(fw, tr); err != nil {
				return fmt.Errorf("Error copying file contents: %v", err), warnings
			}
		case tar.TypeSymlink:
			if !opts.allowSymlinks {
				warnings = append(warnings, fmt.Errorf("Symlinks not allowed: %s", outputPath))
				continue
			}
			// Check if the destination of the symlink is in current directory
			absLinkDest := filepath.Join(absDest, header.Linkname)
			if !Dircheck(absDest, absLinkDest) {
				warnings = append(warnings, fmt.Errorf("Symlink destination outside of current directory: %s -> %s", header.Name, header.Linkname))
				continue
			}
			if opts.validate {
				continue
			}
			// Create symlink
			if err := os.Symlink(header.Linkname, outputPath); err != nil {
				return fmt.Errorf("Error creating symlink: %v", err), warnings
			}
		default:
			warnings = append(warnings, fmt.Errorf("Unable to untar type : %c in file %s", header.Typeflag, header.Name))
			continue
		}
	}
	return nil, warnings
}
