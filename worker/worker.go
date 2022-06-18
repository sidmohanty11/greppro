package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// @property {string} Line - The line of text that matched the pattern
// @property {int} LineNum - The line number of the file where the match was found
// @property {string} Path - The path to the file that was searched
type Result struct {
	Line    string
	LineNum int
	Path    string
}

// `Results` is a struct that contains a slice of `Result`s.
// @property {[]Result} Inner - This is the array of results that we want to return.
type Results struct {
	Inner []Result
}

// NewResult is a function that takes a string, an int, and a string and returns a Result
func NewResult(line string, lineNum int, path string) Result {
	return Result{line, lineNum, path}
}

// It opens a file, scans it line by line, and if the line contains the string we're looking for, it
// adds it to a list of results
func FindInFile(path string, find string) *Results {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	results := Results{make([]Result, 0)}

	scanner := bufio.NewScanner(file)

	lineNum := 1

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), find) {
			r := NewResult(scanner.Text(), lineNum, path)
			results.Inner = append(results.Inner, r)
		}
		lineNum++
	}

	if len(results.Inner) == 0 {
		return nil
	}

	return &results
}
