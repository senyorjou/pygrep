package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/karrick/godirwalk"
)

type semaphore chan bool

func publishPaths(path string, paths chan<- string) {
	godirwalk.Walk(path, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			paths <- osPathname
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})
	close(paths)
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
	const WORKERS = 10
	root := "."
	pathsCh := make(chan string)

	fmt.Printf("Looking at files on %s\n", root)

	re := regexp.MustCompile(`LOL`)

	go publishPaths(root, pathsCh)

	sem := make(semaphore, WORKERS)
	var wg sync.WaitGroup
	i := 0
	for path := range pathsCh {
		sem.Acquire()
		wg.Add(1)
		go func(path string) {
			find(re, path)
			sem.Release()
			wg.Done()
		}(path)
		i += 1
	}
	wg.Wait()

	fmt.Printf("Scanned %d files\n", i)
}
