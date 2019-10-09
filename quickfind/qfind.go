package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var dirs []string
var r bool
var fn string
var verbose bool
var fuzzy bool

var seam = make(chan struct{}, 20)

func main() {
	var wg sync.WaitGroup
	files := make(chan string)
	for _, dir := range dirs {
		wg.Add(1)
		go walkDirs(dir, fn, &wg, files)
	}

	go func() {
		wg.Wait()
		close(files)
	}()

	if !verbose {

	loop:
		for {
			select {
			case dir, ok := <-files:
				if !ok {
					break loop
				}
				println(dir)
			}

		}
	} else {
		println("find files : ")
		for fs := range files {
			println(fs)
		}
	}

}

func walkDirs(dir string, target string, wg *sync.WaitGroup, files chan<- string) {
	defer wg.Done()
	if dir == "" {
		return
	}
	info, err := os.Stat(dir)
	if err != nil {
		println("not a file name or direct : " + dir)
		os.Exit(9)
	}
	if matchFile(info.Name(), fn) {
		//println(info.ModTime().Format("2006-01-02 15:04:05"))
		files <- dir
	}

	if info.IsDir() {
		infos := getChildFiles(dir)
		for _, info := range infos {
			wg.Add(1)
			subDir := filepath.Join(dir, info.Name())
			go walkDirs(subDir, target, wg, files)
		}
	}
}

func getChildFiles(dir string) []os.FileInfo {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("walk dir error : %s %v \n", dir, err)
	}
	return dirs
}

func matchFile(fn string, target string) bool {
	seam <- struct{}{}
	defer func() { <-seam }()
	if r {
		reg := regexp.MustCompile(target)
		return reg.MatchString(fn)

	} else {
		if fuzzy {
			return strings.HasPrefix(fn, target) ||
				strings.Contains(fn, target) ||
				strings.HasSuffix(fn, target)
		}
		return fn == target
	}
}

func init() {

	initParam()

	if fn == "" {
		fmt.Println("specified file name is empty")
		os.Exit(2)
	}

}

func initParam() {
	flag.BoolVar(&r, "r", false, "use regular expression match target file.default false")
	flag.StringVar(&fn, "f", "", "specified file name")
	flag.BoolVar(&verbose, "v", false, "only show search result")
	flag.BoolVar(&fuzzy, "t", false, "fuzzy matching input file name")
	flag.Parse()
	dirs = flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
}
