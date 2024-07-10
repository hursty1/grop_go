package main

import (
	"flag"
	"fmt"
	"grop/config"
	"os"
	"strings"
)

func main() {
	filename := flag.Bool("filename", false, "Print filename")
	ignoreCase := flag.Bool("ignore-case", false, "Ignore case when searching")
	recursive := flag.Bool("recursive", false, "Recursive directory searching")

	//short hands
	flag.BoolVar(ignoreCase, "i", false, "Ignore case when searching")
	flag.BoolVar(filename, "f", false, "Print filename")
	flag.BoolVar(recursive, "r", false, "Recursive directory searching")
	args := os.Args[1:]
	flags := []string{}
	positional := []string{}
	for _, arg := range args {
		// fmt.Println(arg)
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			positional = append(positional, arg)
		}
	}
	// fmt.Println(flags)
	// fmt.Println(positional)
	// flag.Parse() //replaced by line below
	flag.CommandLine.Parse(flags)
	// Positional arguments
	// args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: go run main.go <query> <file> [--ignore-case] [--filename] [--recursive]")
		return
	}
	query := positional[0]
	file := positional[1]
	// fmt.Println(*filename)
	
	configArgs := config.Args{
		Query:      query,
		File:       file,
		IgnoreCase: *ignoreCase,
		Filename:   *filename,
		Recursive:  *recursive,
	}
	var err error
	configSetting, err := config.BuildConfig(configArgs)
	// _, _ = configSetting, err

	if (err != nil) {
		fmt.Println("Error has occured. ", err)
		return
	}

	err = config.Run(configSetting)
	if err != nil{
		fmt.Println("Error has occured. ", err)
	}
	// _ = errMain //debug

	
}
