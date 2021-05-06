package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Lines interface {
	Sum() int
}

type FileLines struct {
	path      string
	filenames []string
	pipe      chan int
}

func filenamesFromPath(path string) []string {
	var filenames []string
	err := filepath.Walk(
		path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if !info.IsDir() {
				filenames = append(filenames, path)
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	return filenames
}

func countLines(filename string) int {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	count := bytes.Count(content, []byte{'\n'})
	fmt.Printf("Lines in '%s': %d\n", filename, count)
	return count
}

func NewFileLines(path string) FileLines {
	filenames := filenamesFromPath(path)
	return FileLines{
		path:      path,
		filenames: filenames,
		pipe:      make(chan int, len(filenames)),
	}
}

func (lines FileLines) Sum() int {
	sum := 0
	for _, filename := range lines.filenames {
		go func(fname string) {
			count := countLines(fname)
			lines.pipe <- count
		}(filename)
	}
	for range lines.filenames {
		sum = sum + <-lines.pipe
	}
	return sum
}

func main() {
	lines := NewFileLines("dataset")
	fmt.Printf("Total of lines: %d\n", lines.Sum())
}
