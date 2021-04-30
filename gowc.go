package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func countLines(wg *sync.WaitGroup, pipe chan int, file string) {
	defer wg.Done()
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	content := string(bytes)
	count := strings.Count(content, "\n")
	fmt.Printf("Lines in '%s': %d\n", file, count)
	pipe <- count
}

func sumLines(pipe chan int) int {
	sum := 0
	for {
		if count, ok := <-pipe; ok {
			sum = sum + count
		} else {
			fmt.Println("saindo!")
			break
		}
	}
	return sum
}

func main() {
	wg := new(sync.WaitGroup)
	pipe := make(chan int, 200)
	var files []string
	root := "dataset"
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		wg.Add(1)
		go countLines(wg, pipe, file)
	}
	sum := sumLines(pipe)
	wg.Wait()
	close(pipe)
	fmt.Printf("Total: %d\n,", sum)
}
