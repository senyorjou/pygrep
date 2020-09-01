package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getPaths(path string) []string {
	var paths []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing a path %q: %v\n", path, err)
			return err
		}
		paths = append(paths, path)
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", ".", err)
	}
	return paths
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

func genPaths(paths []string, ch chan string) {
	for _, path := range paths {
		ch <- path
	}
	close(ch)
}

func find(re *regexp.Regexp, path string) {
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
	paths := getPaths(".")

	fmt.Printf("Looking at %d files\n", len(paths))
	pathsCh := make(chan string)

	var re = regexp.MustCompile(`LOL`)

	go genPaths(paths, pathsCh)
	for path := range pathsCh {
		go find(re, path)
	}

}
