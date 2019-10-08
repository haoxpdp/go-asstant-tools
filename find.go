package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	dirs := flag.Args()
	dir := initParams(dirs)
	println(dir)
	walkDirs(dir,"")
}
func initParams(dirs []string) string {
	var dir string
	if len(dirs) == 0 {
		dir = "."
	} else if len(dirs) == 1 {
		dir = dirs[0]
	} else {
		println("error params")
		os.Exit(1)
	}
	return dir
}

func walkDirs(dir string, target string) {
	info, _ := os.Stat(dir)

	println(dir)

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
	if (err != nil) {
		fmt.Printf("walk dir error : %s %v \n", dir, err)
	}
	return dirs;
}
