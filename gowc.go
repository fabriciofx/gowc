package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func fromDir(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

func filenames(path string) []string {
	var filenames []string
	err := filepath.Walk(path, fromDir(&filenames))
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

func sumLines(pipe chan int) int {
	sum := 0
	numFiles := <-pipe
	for cnt := 0; cnt < numFiles; cnt++ {
		sum = sum + <-pipe
	}
	return sum
}

func main() {
	filenames := filenames("dataset")
	length := len(filenames)
	pipe := make(chan int, length+1)
	pipe <- length
	for _, filename := range filenames {
		go func(fname string) {
			count := countLines(fname)
			pipe <- count
		}(filename)
	}
	total := sumLines(pipe)
	fmt.Printf("Total of lines: %d\n", total)
}
