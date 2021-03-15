package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	_ = readFile()
}
func readFile() error {

	fp := "D:\\tmp\\repo.list"
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	bf := bufio.NewReader(f)
	for {
		line, eir := bf.ReadBytes('\n')

		dir := string(line)
		rmDir(dir)

		if eir != nil {
			if eir == nil {
				return nil
			}
			return eir
		}

	}

}

func rmDir(dir string) {
	existed := true
	dir = strings.Replace(dir, "\n", "", -1)
	index := strings.LastIndex(dir, "\\")
	dir = dir[0:index]
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		existed = false
	}
	if existed {
		fmt.Println("remove : "+dir)
	}
	if existed {
		e := os.RemoveAll(dir)
		if e != nil {
			fmt.Println("delete path [" + dir + "] failed!")
		}
	}
}
