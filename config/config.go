package config

import (
	"bufio"
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
	LineNumber bool
	IsPiped    bool
	Scanner    *bufio.Scanner
	Count	   bool
}
type Config struct {
    Query      string
    FilePath   string
    IgnoreCase bool
    Filename   bool
    Recursive  bool
	LineNumber bool
	IsPiped    bool
	Scanner    *bufio.Scanner
	Count 	   bool
}


func BuildConfig(args Args) (Config, error) {
	var scanner *bufio.Scanner
	if args.Scanner != nil {
		scanner = args.Scanner
	}
	config := Config{
        Query:      args.Query,
        FilePath:   args.File,
        IgnoreCase: args.IgnoreCase,
        Filename:   args.Filename,
        Recursive:  args.Recursive,
		LineNumber: args.LineNumber,
		IsPiped: args.IsPiped,
		Scanner: scanner,
		Count: args.Count,
    }
	// fmt.Println(config)
    return config, nil
}

type Matches struct {
	LineNumber 	int
	FileName 	string
	LineText 	string
	
}

func Run(config Config) (error) {
	var files [] string
	var err error
	var matches []string //slice

	//Data is from piped output 
	if config.IsPiped {
		err := read_Stdin(config)
		if err != nil {
			return err
		}
	}


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
	} else { //matching only to current directory
		matches, err = filepath.Glob(config.FilePath)
		if err != nil {
			return err
		}
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

func read_Stdin(config Config) (error) {
	// fmt.Printf("Scanner is working")
	var lines []byte
	for config.Scanner.Scan() {
		// fmt.Println(config.Scanner.Text())
		lines = append(lines, []byte(config.Scanner.Text()+"\n")...)
	}
	result, err := Search(lines, config.Query, "Piped")
	if err != nil {
		return err
	}
	get_results(result, config)
	return nil
}

func read_results(files []string, config Config) (error) {
	for _, file := range files {
		var result []Matches
		var err error
		contents, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if config.IgnoreCase {
			result, err = CaseinSensitiveSearch(contents, config.Query, file)
		} else {
			result, err = Search(contents, config.Query, file)
		}
		if err != nil{
			return err
		}
		get_results(result, config)
	}
	return nil
}
func get_count_results(results []Matches){
	file_matches := make(map[string]int)

	for _, result := range results {
		file_matches[result.FileName] += 1
	}

	for k, v := range file_matches {
		fmt.Println(k, ": ", v)
	}
}

func get_results(results []Matches, config Config){
	//print count
	if config.Count {
		get_count_results(results)
		
		return
	}

	for _, line := range results {
		// fmt.Println()

		//print count

		var start int
		if config.IgnoreCase{
			start = strings.Index(strings.ToLower(line.LineText), strings.ToLower(config.Query))
		} else {
			start = strings.Index(line.LineText, config.Query)

		}
		
		before := line.LineText[:start]
		match := line.LineText[start : start+len(config.Query)]
		after := line.LineText[start+len(config.Query):]
		if config.Filename{
			fmt.Print(line.FileName + ":")
		}
		if (config.LineNumber){
			// fmt.Print(string(line.LineNumber) + ":")
			fmt.Print(fmt.Sprintf("%v:", line.LineNumber))
		}
		fmt.Print(before + Green + match + Reset + after)
		fmt.Println("")
	}
}

func Search(contents[] byte, query string, filename string) ([]Matches, error) {
	
	var results []Matches
	var count int = 0
	for _, line := range strings.Split(string(contents), "\n") {
		// fmt.Println("line: ", line)
		if strings.Contains(line, query) {
			// line := fmt.Sprintf("%v:%v", count, strings.TrimSpace(line))
			match := Matches{LineNumber: count, FileName: filename, LineText: strings.TrimSpace(line)}
			results = append(results, match)
		}
		count ++ 
	}
	

	return results, nil
}

func CaseinSensitiveSearch(contents[] byte, query string, filename string) ([]Matches, error) { 
	
	var results []Matches
	var count int = 0 //line no
	for _, line := range strings.Split(string(contents), "\n") {
		
		lower_line := strings.ToLower(line)
		lower_query := strings.ToLower(query)

		line := strings.TrimSpace(line) // Temp var to store unedited line
		//assigning each line to preserve the unedited
		
		// fmt.Print("LINE: ", count)
		if strings.Contains(lower_line, lower_query) {

			match := Matches{LineNumber: count, FileName: filename, LineText: strings.TrimSpace(line)} 
			results = append(results, match)
		}
		count ++
	}
	// results.Count = count

	return results, nil
}




