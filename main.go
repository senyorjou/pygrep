package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func getPaths(path string, paths chan<- string) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing a path %q: %v\n", path, err)
			return err
		}
		paths <- path
		return nil
	})

	close(paths)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", ".", err)
	}
}

func getLine(contents []byte, loc []int) string {
	left := loc[0]
	right := loc[1]
	contentLength := len(contents)

	for left > 0 && string(contents[left]) != "\n" && left > (loc[0]-40) {
		left--
	}
	for right <= contentLength && string(contents[right]) != "\n" && right < (loc[1]+40) {
		right++
	}

	return strings.TrimSpace(string(contents[left:right]))

}

func find(re *regexp.Regexp, path string, wg *sync.WaitGroup) {

	defer wg.Done()
	file, _ := os.Stat(path)

	if !file.IsDir() {
		binContent, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		if loc := re.FindIndex(binContent); loc != nil {
			fmt.Printf("> %s ./%s\n", file.Name(), path)
			fmt.Println(getLine(binContent, loc))
		}
	}
}

func main() {
	root := "."
	pathsCh := make(chan string)

	fmt.Printf("Looking at files on %s\n", root)

	re := regexp.MustCompile(`LOL`)

	go getPaths(root, pathsCh)

	var wg sync.WaitGroup
	for path := range pathsCh {
		wg.Add(1)
		go find(re, path, &wg)
	}
	wg.Wait()
}
