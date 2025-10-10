package utilities

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

//depercated
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

//depercated
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
    return config, nil
}

type Matches struct {
	LineNumber 	int
	FileName 	string
	LineText 	string
	
}

func Run(config Args) (error) {
	// var files [] string
	var err error
	// var matches []string //slice

	//Data is from piped output 
	if config.IsPiped {
		err := read_Stdin(config)
		if err != nil {
			return err
		}
	}

	err = crawl_files(config)
	if err != nil {
		return err
	}
	// if config.Recursive {
	// 	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	// 		// Check if the file matches the pattern
	// 		match, err := filepath.Match(config.File, filepath.Base(path))
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if match {
	// 			files = append(files, path)
	// 		}
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// } else { //matching only to current directory
	// 	matches, err = filepath.Glob(config.File)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	
	// files = append(files, matches...)

	// resultErr := read_results(files, config)
	// if resultErr != nil{
	// 	fmt.Println(resultErr)
	// 	return resultErr
	// }
	// fmt.Println(files)
	return nil
}

func crawl_files(config Args) (error) {
	// matches := []string{}
	file_que := make(chan string)
	go func() {
		err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				//  //if not directory send filepath to files slice
				// Check if the file matches the pattern
				match, err := filepath.Match(config.File, filepath.Base(path))
				if err != nil {
					return err
				}
				if match {
					file_que <- path
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Failed to walk %s\n",err)
		}
		close(file_que)
	}()
	var allMatches []Matches
	match_results := make(chan []Matches)
	collectDone := make(chan struct{})
	
	go func() {
		for res := range match_results {
			allMatches = append(allMatches, res...)
		}
		close(collectDone)
	}()

	
	wg := sync.WaitGroup{}
	for f := range file_que{
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			result, err := processFile(file, config)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			match_results <- result
		}(f)
	}
	wg.Wait()
	close(match_results)
	<-collectDone
	//print results
	print_results(allMatches, config)

	return nil
}

func processFile(file_path string, config Args) ([]Matches, error) {
	
	results, err := read_result(file_path, config)
	if err != nil {
		return []Matches{}, err
	}
	return results, nil
}

func read_result(file string, config Args) ([]Matches, error) {
	//Read file and return a slice of Matches
	var result []Matches
	contents, err := os.ReadFile(file)
	if err != nil {
		return result, err
	}
	if config.IgnoreCase {
		result, err = CaseinSensitiveSearch(contents, config.Query, file)
	} else {
		result, err = Search(contents, config.Query, file)
	}
	if err != nil{
		return result, err
	}
	return result, nil
}

func print_results(matches []Matches, config Args) {
	get_results(matches, config)
}
func read_Stdin(config Args) (error) {
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

func read_results(files []string, config Args) (error) {
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

func get_results(results []Matches, config Args){
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





