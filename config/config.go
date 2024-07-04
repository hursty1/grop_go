package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)
type Args struct {
    Query      string
    File       string
    IgnoreCase bool
    Filename   bool
    Recursive  bool
}
type Config struct {
    Query      string
    FilePath   string
    IgnoreCase bool
    Filename   bool
    Recursive  bool
}


func BuildConfig(args Args) (Config, error) {
	config := Config{
        Query:      args.Query,
        FilePath:   args.File,
        IgnoreCase: args.IgnoreCase,
        Filename:   args.Filename,
        Recursive:  args.Recursive,
    }
	// fmt.Println(config)
    return config, nil
}

func Run(config Config) (error) {
	// path := config.FilePath
	var files [] string
	var err error
	var matches []string //slice
	if config.Recursive {
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Check if the file matches the pattern
			match, err := filepath.Match(config.FilePath, filepath.Base(path))
			if err != nil {
				return err
			}
			if match {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		matches, err = filepath.Glob(config.FilePath)

	}
	if err != nil {
		return err
	}
	files = append(files, matches...)

	resultErr := read_results(files, config)
	if resultErr != nil{
		fmt.Println(resultErr)
		return resultErr
	}
	// fmt.Println(files)
	return nil
}

func read_results(files []string, config Config) (error) {
	for _, file := range files {
		var result []string
		var err error
		contents, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if config.IgnoreCase {
			result, err = CaseinSensitiveSearch(contents, config.Query)
		} else {
			result, err = Search(contents, config.Query)
		}
		if err != nil{
			return err
		}
		get_results(result, config, file)

	}


	

	return nil
}

func get_results(results []string, config Config, filename string){
	for _, line := range results {
		fmt.Println()
		var start int
		if config.IgnoreCase{
			start = strings.Index(strings.ToLower(line), strings.ToLower(config.Query))
		} else {
			start = strings.Index(line, config.Query)

		}
		
		before := line[:start]
		match := line[start : start+len(config.Query)]
		after := line[start+len(config.Query):]
		if config.Filename{
			fmt.Print(filename + ": ")
		}
		fmt.Print(before + Green + match + Reset + after)
	}
}

func Search(contents[] byte, query string) ([]string, error) {
	
	var results []string
	
	for _, line := range strings.Split(string(contents), "\n") {
		// fmt.Println("line: ", line)
		if strings.Contains(line, query) {
			results = append(results, strings.TrimSpace(line))
		}
	}
	

	return results, nil
}

func CaseinSensitiveSearch(contents[] byte, query string) ([]string, error) { 
	
	var results []string
	
	for _, line := range strings.Split(string(contents), "\n") {
		// fmt.Println("line: ", line)
		lower_line := strings.ToLower(line)
		lower_query := strings.ToLower(query)
		if strings.Contains(lower_line, lower_query) {
			results = append(results, strings.TrimSpace(line))
		}
	}
	

	return results, nil
}




