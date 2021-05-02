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

func countLines(pipe chan int, filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	count := bytes.Count(content, []byte{'\n'})
	fmt.Printf("Lines in '%s': %d\n", filename, count)
	pipe <- count
}

func sumLines(pipe chan int, numFiles int) int {
	sum := 0
	for cnt := 0; cnt < numFiles; cnt++ {
		sum = sum + <-pipe
	}
	return sum
}

func main() {
	pipe := make(chan int, 200)
	filenames := filenames("dataset")
	for _, filename := range filenames {
		go countLines(pipe, filename)
	}
	total := sumLines(pipe, len(filenames))
	fmt.Printf("Total: %d\n", total)
}
