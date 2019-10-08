package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var dirs []string
var r bool
var fn string
var verbose bool

func main() {
	for _, dir := range dirs {
		walkDirs(dir, fn)
	}
}

func walkDirs(dir string, target string) {
	info, _ := os.Stat(dir)
	if matchFile(info.Name(), fn) {
		println(dir)
	}

	if info.IsDir() {
		infos := getChildFiles(dir)
		for _, info := range infos {
			subDir := filepath.Join(dir, info.Name())
			walkDirs(subDir, target)
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
	if r {
		reg := regexp.MustCompile(target)
		return reg.MatchString(fn)

	} else {
		return strings.HasPrefix(fn, target) ||
			strings.Contains(fn, target) ||
			strings.HasSuffix(fn, target)
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
	flag.BoolVar(&verbose, "v", false, "show search directs")
	flag.Parse()
	dirs = flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
}
