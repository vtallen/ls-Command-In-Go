package main

import (
	"fmt"
	"io/fs"
	"os"
)

func main() {
	fmt.Println("Hello world!")
	args := os.Args
	fmt.Println(args)
	fmt.Println(len(args))

	var hasFlags bool
	var flags string
	flags = "nil"
	if len(args) > 1 {
		flags = args[1]
		if flags[0] == '-' {
			hasFlags = true
			flags = flags[1:]
		}
	}

	var path string
	var err error
	if len(args) < 3 {
		path, err = os.Getwd()
		if err != nil {
			fmt.Println("Erorr getting current directory")
			os.Exit(1)
		}
	} else {
		if hasFlags {
			// handle flag parsing
			path = args[2]
		} else {
			path = args[1]
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		os.Exit(1)
	}

	filesInfo := make([]fs.FileInfo, 0, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		// fmt.Printf("type %T", info)
		filesInfo = append(filesInfo, info)
	}

	fmt.Println()
}
