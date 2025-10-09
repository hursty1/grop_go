package utilities

import "strings"

func Search(contents []byte, query string, filename string) ([]Matches, error) {

	var results []Matches
	var count int = 0
	for _, line := range strings.Split(string(contents), "\n") {
		// fmt.Println("line: ", line)
		if strings.Contains(line, query) {
			// line := fmt.Sprintf("%v:%v", count, strings.TrimSpace(line))
			match := Matches{LineNumber: count, FileName: filename, LineText: strings.TrimSpace(line)}
			results = append(results, match)
		}
		count++
	}

	return results, nil
}

func CaseinSensitiveSearch(contents []byte, query string, filename string) ([]Matches, error) {

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
		count++
	}
	// results.Count = count

	return results, nil
}
