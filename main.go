package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
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

type Flags struct {
	LongListing   *bool
	HumanReadable *bool
	Recursive     *bool
	SortTime      *bool
	SortSize      *bool
	Reverse       *bool
	NoColors      *bool
	ShowHidden    *bool
	ShowINodes    *bool
	Path          string
}

func main() {
	FlagsBool := ParseArgs()
	DebugArgs(FlagsBool)

	// args := os.Args

	// var hasFlags bool
	// var flags string
	// flags = "nil"
	// if len(args) > 1 {
	// 	flags = args[1]
	// 	if flags[0] == '-' {
	// 		hasFlags = true
	// 		flags = flags[1:]
	// 	}
	// }

	// var path string
	// var err error
	// if len(args) < 3 {
	// 	path, err = os.Getwd()
	// 	if err != nil {
	// 		fmt.Println("Erorr getting current directory")
	// 		os.Exit(1)
	// 	}
	// } else {
	// 	if hasFlags {
	// 		// handle flag parsing
	// 		path = args[2]
	// 	} else {
	// 		path = args[1]
	// 	}
	// }

	// filesInfo := filterHidden(getFilesInfo(path))
	// for _, info := range filesInfo {
	// 	fmt.Printf("%s ", getColorFilename(info))
	// }

	// fmt.Println()
}

func ParseArgs() *Flags {
	var FlagsBool Flags
	// Define flags
	FlagsBool.LongListing = flag.Bool("l", false, "Use long listing format")
	FlagsBool.HumanReadable = flag.Bool("h", false, "Print sizes in human readable format")
	FlagsBool.Recursive = flag.Bool("R", false, "List subdirectories recursively")
	FlagsBool.SortTime = flag.Bool("t", false, "Sort by modification time")
	FlagsBool.SortSize = flag.Bool("S", false, "Sort by file size")
	FlagsBool.Reverse = flag.Bool("r", false, "Reverse the order of sort")
	FlagsBool.NoColors = flag.Bool("G", false, "Disable colorized output")

	// Define flags related to filtering
	FlagsBool.ShowHidden = flag.Bool("a", false, "Show hidden files")
	FlagsBool.ShowINodes = flag.Bool("i", false, "Print the index number of each file")

	argv := os.Args[1:]
	argc := len(os.Args[1:])

	var err error

	// Case when vls
	if argc == 0 {
		FlagsBool.Path, err = os.Getwd()

	} else if argc == 1 {
		/*
		   Case when:
		   vls <path>
		   vls -l
		   vls -lah
		*/

		// Catches multiple argument case
		if argv[0][0] == '-' && len(argv[0]) > 2 {
			ParseMultiFlags()
			FlagsBool.Path, err = os.Getwd()
		} else if argv[0][0] == '-' { // Catches the single flag case
			flag.Parse()
			FlagsBool.Path, err = os.Getwd()
		} else { // Catches the case when just a path is given
			FlagsBool.Path = argv[0]
		}

	} else if argc == 2 { // Catches the case when command is in form vls -lah <path>
		ParseMultiFlags()
		FlagsBool.Path = argv[argc-1]
	} else { // Catches the case in which more than 1 flag is given seperately
		flag.Parse()
		leftover := flag.Args()

		if len(leftover) > 1 {
			flag.CommandLine.Usage()
		}

		FlagsBool.Path = leftover[0]
	}

	if err != nil {
		fmt.Printf("Error parsing args:%s\n", err)
		os.Exit(1)
	}

	return &FlagsBool

}

func DebugArgs(FlagsBool *Flags) {
	fmt.Println("Path:", FlagsBool.Path)
	fmt.Println()
	// Parse command-line arguments
	// flag.Parse()

	// Access the values of the flags
	fmt.Println("Formatting flags:")
	fmt.Println("-l:", *FlagsBool.LongListing)
	fmt.Println("-h:", *FlagsBool.HumanReadable)
	fmt.Println("-R:", *FlagsBool.Recursive)
	fmt.Println("-t:", *FlagsBool.SortTime)
	fmt.Println("-S:", *FlagsBool.SortSize)
	fmt.Println("-r:", *FlagsBool.Reverse)
	fmt.Println("-G:", *FlagsBool.NoColors)

	fmt.Println("\nFiltering flags:")
	fmt.Println("-a:", *FlagsBool.ShowHidden)
	fmt.Println("-i:", *FlagsBool.ShowINodes)

	// Access non-flag arguments (if any)
	fmt.Println("\nNon-flag arguments:")
	fmt.Println(flag.Args())

}

func ParseMultiFlags() {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-") {
			for _, flagChar := range arg[1:] {
				flag := flag.Lookup(string(flagChar))
				if flag == nil {
					fmt.Printf("Unknown flag: -%c\n", flagChar)
					os.Exit(1)
				}
				flag.Value.Set("true")
			}
		} else {
			break
		}
	}
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

func filterHidden(filesInfo []fs.FileInfo) []fs.FileInfo {
	noHidden := make([]fs.FileInfo, 0, 0)
	for _, file := range filesInfo {
		if file.Name()[0] != '.' {
			noHidden = append(noHidden, file)
		}
	}

	return noHidden
}

func GetColorFilename(fileinfo fs.FileInfo) string {
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
