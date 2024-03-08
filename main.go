package main

import (
	"fmt"
	"io/fs"
	"os"
)

const (
	RESET  = "\033[0m"
	BOLD   = "\033[1m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PURPLE = "\033[35m"
	CYAN   = "\033[36m"
	GREY   = "\033[37m"
)

func main() {
	args := os.Args

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

	filesInfo := getFilesInfo(path)
	for _, info := range filesInfo {
		fmt.Printf("%s ", getColorFilename(info))
	}

	fmt.Println()
}

func filterHidden(filesInfo []fs.FileInfo) []fs.FileInfo {
	return filesInfo
}

func getFilesInfo(path string) []os.FileInfo {
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
	return filesInfo
}

func getColorFilename(fileinfo fs.FileInfo) string {
	var color string
	if fileinfo.IsDir() {
		color = BOLD + BLUE
	} else {
		if fileinfo.Mode()&os.ModePerm&0100 != 0 { // Tests if the file is an executable
			color = BOLD + GREEN
		} else {
			color = ""
		}
	}

	return color + fileinfo.Name() + RESET
}
