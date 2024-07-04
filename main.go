package main

import (
	"flag"
	"fmt"
	"grop/config"
)

func main() {
	filename := flag.Bool("filename", false, "Print filename")
	ignoreCase := flag.Bool("ignore-case", false, "Ignore case when searching")
	recursive := flag.Bool("recursive", false, "Recursive directory searching")

	//short hands
	flag.BoolVar(ignoreCase, "i", false, "Ignore case when searching")
	flag.BoolVar(filename, "f", false, "Print filename")
	flag.BoolVar(recursive, "r", false, "Recursive directory searching")

	flag.Parse()

	// Positional arguments
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: go run main.go <query> <file> [--ignore-case] [--filename] [--recursive]")
		return
	}
	query := args[0]
	file := args[1]
	// fmt.Println(*ignoreCase)
	
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
