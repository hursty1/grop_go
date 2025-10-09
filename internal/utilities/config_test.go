package utilities

import (
	"bufio"
	"strings"
	"testing"
)


func TestReadStdin(t *testing.T) {
	input := "TESTING \n new console output"
	scanner := bufio.NewScanner(strings.NewReader(input))
	config := Args{
		Query:     "new",
		File:   "",
		IgnoreCase: false,
		Filename:   false,
		Recursive:  false,
		LineNumber: false,
		IsPiped:    true,
		Scanner:    scanner,
		Count: false,
}
	err := read_Stdin(config)
	if err != nil {
		t.Errorf("read_Stdin() failed: %v", err)
	}
}


