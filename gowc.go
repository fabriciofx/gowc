package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Files struct {
	path  string
	names []string
	pipe  chan int
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

func NewFiles(path string) Files {
	names := filenamesFromPath(path)
	return Files{
		path:  path,
		names: names,
		pipe:  make(chan int, len(names)),
	}
}

func (files Files) FilesLinesSum() int {
	sum := 0
	for _, filename := range files.names {
		go func(fname string) {
			count := countLines(fname)
			files.pipe <- count
		}(filename)
	}
	for range files.names {
		sum = sum + <-files.pipe
	}
	return sum
}

func main() {
	files := NewFiles("dataset")
	fmt.Printf("Total of lines: %d\n", files.FilesLinesSum())
}
