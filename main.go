package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sort"
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
	ArgsFlags := ParseArgs()
	// DebugArgs(ArgsFlags)
	// fmt.Println()

	filesInfo := GetFilesInfo(ArgsFlags.Path)
	if *ArgsFlags.LongListing {
		PrintLongListing(ArgsFlags, filesInfo)
	} else {
		PrintNormalListing(ArgsFlags, filesInfo)
	}
}

func PrintUsage() {
	const USAGE string = "A copy of the ls command written in go\n Examples:\n\tvls <path>\n\tvls -lah <path>\n\tvls -l -a -h <path>\n\tvls -lah\n\tvls -l -a -h"
	fmt.Println(USAGE)
	flag.PrintDefaults()
}

func ParseArgs() *Flags {
	var ArgsFlags Flags
	// Define flags
	ArgsFlags.LongListing = flag.Bool("l", false, "Use long listing format")
	ArgsFlags.HumanReadable = flag.Bool("h", false, "Print sizes in human readable format")
	ArgsFlags.Recursive = flag.Bool("R", false, "List subdirectories recursively")
	ArgsFlags.SortTime = flag.Bool("t", false, "Sort by modification time")
	ArgsFlags.SortSize = flag.Bool("S", false, "Sort by file size")
	ArgsFlags.Reverse = flag.Bool("r", false, "Reverse the order of sort")
	ArgsFlags.NoColors = flag.Bool("G", false, "Disable colorized output")

	// Define flags related to filtering
	ArgsFlags.ShowHidden = flag.Bool("a", false, "Show hidden files")
	ArgsFlags.ShowINodes = flag.Bool("i", false, "Print the index number of each file")

	argv := os.Args[1:]
	argc := len(os.Args[1:])

	var err error

	// Case when vls
	if argc == 0 {
		ArgsFlags.Path, err = os.Getwd()

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
			ArgsFlags.Path, err = os.Getwd()
		} else if argv[0][0] == '-' { // Catches the single flag case
			flag.Parse()
			ArgsFlags.Path, err = os.Getwd()
		} else { // Catches the case when just a path is given
			ArgsFlags.Path = argv[0]
		}

	} else if argc == 2 { // Catches the case when command is in form vls -lah <path>
		ParseMultiFlags()
		ArgsFlags.Path = argv[argc-1]
	} else { // Catches the case in which more than 1 flag is given seperately
		flag.Parse()
		leftover := flag.Args()

		if len(leftover) > 1 {
			PrintUsage()
			os.Exit(1)
		}

		ArgsFlags.Path = leftover[0]
	}

	if err != nil {
		fmt.Printf("Error parsing args:%s\n", err)
		PrintUsage()
		os.Exit(1)
	}

	return &ArgsFlags

}

func DebugArgs(ArgsFlags *Flags) {
	fmt.Println("Path:", ArgsFlags.Path)
	fmt.Println()
	// Parse command-line arguments
	// flag.Parse()

	// Access the values of the flags
	fmt.Println("Formatting flags:")
	fmt.Println("-l:", *ArgsFlags.LongListing)
	fmt.Println("-h:", *ArgsFlags.HumanReadable)
	fmt.Println("-R:", *ArgsFlags.Recursive)
	fmt.Println("-t:", *ArgsFlags.SortTime)
	fmt.Println("-S:", *ArgsFlags.SortSize)
	fmt.Println("-r:", *ArgsFlags.Reverse)
	fmt.Println("-G:", *ArgsFlags.NoColors)

	fmt.Println("\nFiltering flags:")
	fmt.Println("-a:", *ArgsFlags.ShowHidden)
	fmt.Println("-i:", *ArgsFlags.ShowINodes)

	// Access non-flag arguments (if any)
	fmt.Println("\nNon-flag arguments:")
	fmt.Println(flag.Args())

}

func ParseMultiFlags() {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-") {
			for _, flagChar := range arg[1:] {
				curFlag := flag.Lookup(string(flagChar))
				if curFlag == nil {
					PrintUsage()
					os.Exit(1)
				}
				curFlag.Value.Set("true")
			}
		} else {
			break
		}
	}
}

func GetFilesInfo(path string) []os.FileInfo {
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

func filterHidden(filesInfo []fs.FileInfo) []fs.FileInfo {
	noHidden := make([]fs.FileInfo, 0, 0)
	for _, file := range filesInfo {
		if file.Name()[0] != '.' {
			noHidden = append(noHidden, file)
		}
	}

	return noHidden
}

func SortName(ArgsFlags *Flags, filesInfo []fs.FileInfo) {
	sort.Slice(filesInfo, func(idxa, idxb int) bool {
		if *ArgsFlags.Reverse {
			return filesInfo[idxa].Name() > filesInfo[idxb].Name()
		} else {
			return filesInfo[idxa].Name() < filesInfo[idxb].Name()
		}
	})
}

func SortSize(ArgsFlags *Flags, filesInfo []fs.FileInfo) {
	sort.Slice(filesInfo, func(idxa, idxb int) bool {
		if *ArgsFlags.Reverse {
			return filesInfo[idxa].Size() < filesInfo[idxb].Size()
		} else {
			return filesInfo[idxa].Size() > filesInfo[idxb].Size()
		}
	})
}

func PrintNormalListing(ArgsFlags *Flags, filesInfo []fs.FileInfo) {
	if *ArgsFlags.Reverse && !*ArgsFlags.SortSize && !*ArgsFlags.SortTime {
		SortName(ArgsFlags, filesInfo)
	} else if *ArgsFlags.SortSize && !*ArgsFlags.SortTime {
		SortSize(ArgsFlags, filesInfo)
	} else {
		SortName(ArgsFlags, filesInfo)
	}

	for _, info := range filesInfo {
		var finalOut string
		if *ArgsFlags.NoColors {
			finalOut = finalOut + info.Name()
		} else {
			finalOut = finalOut + GetColorFilename(info)
		}

		fmt.Printf(finalOut + " ")
	}
	fmt.Println()
}

func PrintLongListing(ArgsFlags *Flags, filesInfo []fs.FileInfo) {

}
