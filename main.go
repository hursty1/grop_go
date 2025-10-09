package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hursty1/grop_go/internal"
)

func main() {
	filename := flag.Bool("filename", false, "Print filename")
	ignoreCase := flag.Bool("ignore-case", false, "Ignore case when searching")
	recursive := flag.Bool("recursive", false, "Recursive directory searching")
	linenumber := flag.Bool("linenumber", false, "Print the line Number")
	count := flag.Bool("count", false, "Prints the count of matches in each file")
	//short hands
	flag.BoolVar(ignoreCase, "i", false, "Ignore case when searching")
	flag.BoolVar(filename, "f", false, "Print filename")
	flag.BoolVar(recursive, "r", false, "Recursive directory searching")
	flag.BoolVar(linenumber, "l", false, "Show Line Number")
	flag.BoolVar(count, "c", false, "Prints the number of matches in each file")


	args := os.Args[1:]
	flags := []string{}
	positional := []string{}
	for _, arg := range args {
		// fmt.Println(arg) //debug
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			positional = append(positional, arg)
		}
	}
	
	fi, _ := os.Stdin.Stat()
	isPiped := (fi.Mode() & os.ModeCharDevice) == 0

	flag.CommandLine.Parse(flags)

	var scanner *bufio.Scanner

	query := ""
	file := ""
	if isPiped {
		scanner = bufio.NewScanner(os.Stdin)
		query = positional[0]
	} else if len(args) < 2 {
		fmt.Println("Usage: go run main.go <query> <file> [--ignore-case] [--filename] [--recursive] [--count] [--linenumber]")
		return
	} else { // assuming 2 positional
		query = positional[0]
		file = positional[1]
	}
	
	configArgs := internal.Args{
		Query:      query,
		File:       file,
		IgnoreCase: *ignoreCase,
		Filename:   *filename,
		Recursive:  *recursive,
		LineNumber: *linenumber,
		IsPiped: isPiped,
		Scanner: scanner,
		Count: *count,
	}
	var err error
	configSetting, err := internal.BuildConfig(configArgs)
	if (err != nil) {
		fmt.Println("Error Building Configuration.", err)
		return	
	}

	err = internal.Run(configSetting)
	if err != nil{
		fmt.Println("Error Executing Specified build config", err)
		return
	}
	
}
