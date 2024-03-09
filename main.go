package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"sort"
	"strings"
	"syscall"
)

// Terminal color codes
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

// Holds the command line arguments, used by printing functions
// to determine how to print output
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

/*********************************************************************************************
*                                                                                            *
* Name: PrintUsage                                                                           *
*                                                                                            *
* Description: Prints the command examples as wells as flags.PrintDefaults() for flags usage *
*                                                                                            *
* Parameters: none                                                                           *
*                                                                                            *
* return: none                                                                               *
**********************************************************************************************/
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
		// TODO: Need to catch the case where the flag is a single flag
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

func FilterHidden(filesInfo []fs.FileInfo) []fs.FileInfo {
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

func SortTime(ArgsFlags *Flags, filesInfo []fs.FileInfo) {
	sort.Slice(filesInfo, func(idxa, idxb int) bool {

		var comp int = filesInfo[idxa].ModTime().Compare(filesInfo[idxb].ModTime())
		if *ArgsFlags.Reverse {
			if comp == -1 || comp == 0 {
				return true
			} else {
				return false
			}
		} else {
			if comp == -1 || comp == 0 {
				return false
			} else {
				return true
			}
		}
	})
}

func GetINode(fileInfo *fs.FileInfo) (uint64, error) {
	if fileInfo == nil {
		panic("Cannot pass nil into GetINode")
	}

	stat, ok := (*fileInfo).Sys().(*syscall.Stat_t)

	if !ok {
		return 0, fmt.Errorf("unable to get inode for file %s\n", (*fileInfo).Name())
	}

	return stat.Ino, nil
}

func GetFilePerms(fileInfo *fs.FileInfo) string {
	mode := (*fileInfo).Mode()

	permissions := "-"

	// Owner permissions
	if mode&0400 != 0 {
		permissions += "r"
	} else {
		permissions += "-"
	}
	if mode&0200 != 0 {
		permissions += "w"
	} else {
		permissions += "-"
	}
	if mode&0100 != 0 {
		permissions += "x"
	} else {
		permissions += "-"
	}

	// Group permissions
	if mode&0040 != 0 {
		permissions += "r"
	} else {
		permissions += "-"
	}
	if mode&0020 != 0 {
		permissions += "w"
	} else {
		permissions += "-"
	}
	if mode&0010 != 0 {
		permissions += "x"
	} else {
		permissions += "-"
	}

	// Other permissions
	if mode&0004 != 0 {
		permissions += "r"
	} else {
		permissions += "-"
	}
	if mode&0002 != 0 {
		permissions += "w"
	} else {
		permissions += "-"
	}
	if mode&0001 != 0 {
		permissions += "x"
	} else {
		permissions += "-"
	}

	return permissions
}

func GetReadableSize(size int64) string {
	if size < 1024 { // Size is measureable in bytes
		return fmt.Sprint(size)
	} else if size >= 1024 && size < 1048576 { // Size is measureable in kilobytes
		return fmt.Sprintf("%.1f", float64(size)/1024.0) + fmt.Sprint("K")
	} else if size >= 1048576 && size < 1073741824 { // Size is measureable in megabytes
		return fmt.Sprintf("%.1f", float64(size)/1048576.0) + fmt.Sprint("M")
	} else { // Size is measureable in terabytes
		return fmt.Sprintf("%.1f", float64(size)/1073741824.0) + fmt.Sprint("T")
	}
}

func SortFilterOnFlags(ArgsFlags *Flags, filesInfo []fs.FileInfo) {
	// Determine how to sort the entries based on the arguments
	if *ArgsFlags.Reverse && !*ArgsFlags.SortSize && !*ArgsFlags.SortTime {
		SortName(ArgsFlags, filesInfo)
	} else if *ArgsFlags.SortSize && !*ArgsFlags.SortTime {
		SortSize(ArgsFlags, filesInfo)
	} else if !*ArgsFlags.SortSize && *ArgsFlags.SortTime {
		SortTime(ArgsFlags, filesInfo)
	} else {
		SortName(ArgsFlags, filesInfo)
	}

	// If -a is not present in args, take out all hidden files from output
	if !*ArgsFlags.ShowHidden {
		filesInfo = FilterHidden(filesInfo)
	}
}

func PrintNormalListing(ArgsFlags *Flags, filesInfo []fs.FileInfo) {
	// Uses the argument flags to sort and filter the output
	SortFilterOnFlags(ArgsFlags, filesInfo)

	for _, info := range filesInfo {
		var finalOut string

		if *ArgsFlags.ShowINodes {
			inode, ok := GetINode(&info)
			if ok != nil {
				fmt.Printf("%s", ok)
			}

			finalOut = finalOut + fmt.Sprint(inode) + " "
		}

		if *ArgsFlags.NoColors {
			finalOut = finalOut + info.Name()
		} else {
			finalOut = finalOut + GetColorFilename(info)
		}

		fmt.Printf(finalOut + " ")
	}
	if len(filesInfo) > 0 {
		fmt.Println()
	}
}

func PrintTable(table [][]string) {
	cols := len(table[0])
	colSizes := make([]int, cols)

	// Calculate the column sizes
	for _, row := range table {
		for coli, col := range row {
			length := len(col)
			if length > colSizes[coli] {
				colSizes[coli] = length
			}
		}
	}

	// Print out the table
	for _, row := range table {
		var outRow string

		for coli, col := range row {
			outRow = outRow + fmt.Sprintf("%-*s ", colSizes[coli], col)
		}

		outRow = strings.TrimSpace(outRow)
		fmt.Print(outRow + "\n")
	}
}

func PrintLongListing(ArgsFlags *Flags, filesInfo []fs.FileInfo) {
	SortFilterOnFlags(ArgsFlags, filesInfo)

	// Allocate the memory that will store the info for each file
	outTable := make([][]string, len(filesInfo))
	for idx := 0; idx < len(filesInfo); idx++ {
		outTable[idx] = make([]string, 0)
	}

	var totalSize int64
	for idx, info := range filesInfo {

		stat, syscallOk := info.Sys().(*syscall.Stat_t)
		if !syscallOk {
			fmt.Printf("syscall Stat_t failed: %v\n", syscallOk)
		}

		// File inode
		var inode string
		if *ArgsFlags.ShowINodes {
			inodeInt, ok := GetINode(&info)
			if ok != nil {
				fmt.Printf("%s", ok)
			}
			inode = fmt.Sprint(inodeInt)
		}
		outTable[idx] = append(outTable[idx], inode)

		// File permissions
		var permissions string = GetFilePerms(&info)
		outTable[idx] = append(outTable[idx], permissions)

		// Number of hard links
		var numLinks string = fmt.Sprint(stat.Nlink)
		outTable[idx] = append(outTable[idx], numLinks)

		// Get the owner of the file
		ownerUsr, ok := user.LookupId(fmt.Sprint(stat.Uid))
		if ok != nil {
			fmt.Printf("%s\n", ok)
		}
		var owner string = ownerUsr.Username
		outTable[idx] = append(outTable[idx], owner)

		// Group of the file
		var group string
		groupUsr, ok := user.LookupGroupId(fmt.Sprint(stat.Gid))
		if ok != nil {
			fmt.Printf("%s\n", ok)
		}
		group = groupUsr.Name
		outTable[idx] = append(outTable[idx], group)

		// Size of the file
		var size string
		if !*ArgsFlags.HumanReadable {
			size = fmt.Sprint(info.Size())
		} else {
			size = GetReadableSize(info.Size())
		}
		totalSize = totalSize + int64(info.Size())
		outTable[idx] = append(outTable[idx], size)

		// Date/time modified
		var dateTime string = info.ModTime().Format("Jan 02 15:04")
		outTable[idx] = append(outTable[idx], dateTime)

		// Get the filename
		var filename string
		if *ArgsFlags.NoColors {
			filename = info.Name()
		} else {
			filename = GetColorFilename(info)
		}
		outTable[idx] = append(outTable[idx], filename)
	}
	if *ArgsFlags.HumanReadable {
		fmt.Printf("total %s\n", GetReadableSize(totalSize))
	} else {
		fmt.Printf("total %v\n", totalSize)
	}
	PrintTable(outTable)
}
