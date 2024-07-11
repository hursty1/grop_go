package main

import (
	"grop/config"
	"testing"
)

func TestSearch(t *testing.T) {
	contents := `\
GO:
safe, fast, productive.
Pick three.
Duct tape.
	`
	query := "fast"
	expected := []string{"safe, fast, productive."}
	result, err := config.Search([]byte(contents), query, "Filename")
	if err != nil {
		t.Fatalf("Expected no Error, got %v", err)
	}
	resultStr := []string{result[0].LineText}
	if !equal(resultStr, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
func TestCaseSearch(t *testing.T) {
	contents := `\
GO:
safe, fast, productive.
Pick three.
Duct tape.
	`
	query := "Fast"
	expected := []string{}
	result, err := config.CaseinSensitiveSearch([]byte(contents), query, "Filename")
	if err != nil {
		t.Fatalf("Expected no Error, got %v", err)
	}
	resultStr := []string{result[0].LineText}
	if equal(resultStr, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
func TestCaseQuerySearch(t *testing.T) {
	contents := `\
GO:
safe, fast, productive.
Pick three.
Duct tape.
	`
	query := "Fast"
	expected := []string{"safe, fast, productive."}
	result, err := config.CaseinSensitiveSearch([]byte(contents), query, "FILENAME")
	if err != nil {
		t.Fatalf("Expected no Error, got %v", err)
	}
	resultStr := []string{result[0].LineText}
	if !equal(resultStr, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
func TestNoResultSearch(t *testing.T) {
	contents := `\
GO:
safe, fast, productive.
Pick three.
Duct tape.
	`
	query := "Rust"
	expected := []string{}
	result, err := config.Search([]byte(contents), query, "filename")
	
	if err != nil {
		t.Fatalf("Expected no Error, got %v", err)
	}
	resultStr := []string{result[0].LineText}
	if !equal(resultStr, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}