package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func visit(files *[]string) filepath.WalkFunc {
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

func countLines(pipe chan int, file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	count := bytes.Count(content, []byte{'\n'})
	fmt.Printf("Lines in '%s': %d\n", file, count)
	pipe <- count
}

func sumLines(pipe chan int) int {
	sum := 0
	for {
		count, ok := <-pipe
		if ok {
			sum = sum + count
		} else {
			break
		}
	}
	return sum
}

func main() {
	pipe := make(chan int, 200)
	var files []string
	root := "dataset"
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		go countLines(pipe, file)
	}
	sum := sumLines(pipe)
	fmt.Printf("Total: %d\n", sum)
}
