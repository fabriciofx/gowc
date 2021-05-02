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
	pipe := make(chan int, 200)
	filenames := filenames("dataset")
	pipe <- len(filenames)
	for _, filename := range filenames {
		go func(fname string) {
			count := countLines(fname)
			pipe <- count
		}(filename)
	}
	total := sumLines(pipe)
	fmt.Printf("Total: %d\n", total)
}
