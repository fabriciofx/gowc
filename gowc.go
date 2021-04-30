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

type linesum struct {
	file  string
	count int
}

func countLines(files []string) chan linesum {
	ch := make(chan linesum, 200)

	go func() {
		for _, f := range files {
			b, err := os.ReadFile(f)
			if err != nil {
				panic(err)
			}
			count := bytes.Count(b, []byte{'\n'})
			ch <- linesum{f, count}
		}
		close(ch)
	}()
	return ch
}

func main() {
	var files []string
	root := "dataset"
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	linech := countLines(files)
	sum := 0
	for l := range linech {
		sum += l.count
		fmt.Printf("Lines in '%s': %d\n", l.file, l.count)
	}

	fmt.Printf("Total: %d\n,", sum)
}
