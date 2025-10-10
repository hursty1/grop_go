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

type Matches struct {
	LineNumber 	int
	FileName 	string
	LineText 	string
	
}

func Run(config Args) (error) {
	var err error
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

	return nil
}

func crawl_files(config Args) (error) {

	file_que := make(chan string)
	go func() {
		err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {//if not directory send filepath to files slice
				
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
	
	done := make(chan struct{})
	go func() {
		
		for match := range match_results {
			print_results(match, config)
		}
		close(done)
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
	<- done

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
			fmt.Printf("%v:\t", line.LineNumber)
		}
		
		fmt.Print(before + Green + match + Reset + after)
		fmt.Println("")
	}
}