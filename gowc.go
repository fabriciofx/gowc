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

func countLines(filenames []string) int {
	pipe := make(chan int)
	for _, filename := range filenames {
		go func(f string) {
			content, err := os.ReadFile(f)
			if err != nil {
				panic(err)
			}
			count := bytes.Count(content, []byte{'\n'})
			fmt.Printf("Lines in '%s': %d\n", f, count)
			pipe <- count
		}(filename)
	}
	sum := 0
	for range filenames {
		count, ok := <-pipe
		if ok {
			sum = sum + count
		}
	}
	return sum
}

func main() {
	filenames := filenames("dataset")
	fmt.Printf("Total: %d\n", countLines(filenames))
}
