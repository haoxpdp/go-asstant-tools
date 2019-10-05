package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var verbosse = flag.Bool("v", false, "show verbose progress message")

var sema = make(chan struct{}, 20)

var rootMapSize = make(map[string]int64)
var rootMapNum = make(map[string]int64)

type fileSize struct {
	Path string
	Size int64
}

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	fileSizes := make(chan int64)
	filePathSize := make(chan fileSize)
	var wg sync.WaitGroup

	if len(roots) == 1 {
		tmpRoots := []string{}
		absDir, _ := filepath.Abs(roots[0])
		for _, dir := range dirents(absDir) {
			tmpRoots = append(tmpRoots, filepath.Join(absDir, dir.Name()))
		}
		roots = tmpRoots
	}

	for _, dir := range roots {
		absPath, _ := filepath.Abs(dir)
		wg.Add(1)
		rootMapSize[absPath] = 0
		go walkDir(absPath, &wg, filePathSize, fileSizes)
	}
	go func() {
		wg.Wait()
		close(fileSizes)
		close(filePathSize)
	}()

	var fileNum, totalSize int64

	// 定期输出结果
	var ticket <-chan time.Time
	if *verbosse {
		ticket = time.Tick(500 * time.Millisecond)
	}
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop
			}
			fileNum++
			totalSize += size
		case <-ticket:
			printDiskUsages(fileNum, totalSize)
		case fileSize := <-filePathSize:
			addRootPathSize(rootMapSize,fileSize)
		}

	}

	for k, v := range rootMapSize {
		fmt.Printf("%s \t\t\t %.1f \n", k, float64(v)/1e9)
	}

	fmt.Println("total: ")
	printDiskUsages(fileNum, totalSize)
}

func recovery() {
	if r := recover(); r != nil {
		fmt.Println("recovered:", r)
	}
}

func printDiskUsages(fileNum, fileSize int64) {
	fmt.Printf("%d files, %.1f gb\n", fileNum, float64(fileSize)/1e9)
}

func walkDir(dir string, wg *sync.WaitGroup, filePathSize chan<- fileSize, filesizes chan<- int64) {
	defer recovery()
	defer wg.Done()
	fileInfo, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("du1 : %v \n", err)
	}
	if fileInfo.IsDir() {
		for _, entry := range dirents(dir) {
			if entry == nil {
				return
			}
			wg.Add(1)
			subDir := filepath.Join(dir, entry.Name())
			walkDir(subDir, wg, filePathSize, filesizes)
		}
	} else {
		tmpFileSize := fileSize{Path: dir, Size: fileInfo.Size()}
		filePathSize <- tmpFileSize
		filesizes <- fileInfo.Size()
	}
}

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("du1 : %v \n", err)
		return nil
	}
	return entries
}

func addRootPathSize(rootMap map[string]int64, file fileSize) {
	for k, _ := range rootMap {
		if strings.HasPrefix(file.Path, k) {
			rootMap[k] += file.Size
			break
		}
	}
}
