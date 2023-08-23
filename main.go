package main

import (
	"flag"
	"fmt"
	"strings"
)

const (
	TAR_GZ Filetype = iota
	ZIP    Filetype = iota
	OTHER  Filetype = iota
)

type Filetype int

type Options struct {
	maxSize       int64
	allowSymlinks bool
	validate      bool
}

func filetype(src string) Filetype {
	if strings.HasSuffix(src, ".tar.gz") {
		return TAR_GZ
	} else if strings.HasSuffix(src, ".zip") {
		return ZIP
	} else {
		return OTHER
	}
}

func main() {
	var srcFlag = flag.String("src", "", "Source file")
	var destFlag = flag.String("dest", "", "Destination directory")
	var maxSizeFlag = flag.Int64("maxSize", 100000000, "Maximum size of extracted file")
	var allowSymlinksFlag = flag.Bool("allowSymlinks", false, "Allow symlinks")
	var validate = flag.Bool("validate", false, "Validate archive without extracting")
	flag.Parse()
	if *srcFlag == "" || *destFlag == "" {
		flag.PrintDefaults()
		return
	}
	options := Options{maxSize: *maxSizeFlag, allowSymlinks: *allowSymlinksFlag, validate: *validate}
	// Check type of file, and switch accordingly
	// if file is tar.gz, untar
	// if file is zip, unzip
	// if file is neither, return error
	var err error
	var warnings []error
	switch filetype(*srcFlag) {
	case TAR_GZ:
		err, warnings = Untar(*srcFlag, *destFlag, options)
	case ZIP:
		err, warnings = Unzip(*srcFlag, *destFlag, options)
	case OTHER:
		fmt.Println("File type not supported")
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	if len(warnings) > 0 {
		fmt.Println("Errors encountered while untarring")
		for _, err := range warnings {
			fmt.Println(err)
		}
	} else if *validate {
		fmt.Println("Archive is valid")
	} else {
		fmt.Println("Archive extracted successfully")
	}
}
