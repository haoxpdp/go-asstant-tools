package main

import (
	"fmt"
	"os"
)

func main() {
	dir := "D:\\mvn-repo\\com\\yonghui\\super-species-process-api\\1.0-SNAPSHOT\\super-species-process-api-1.0-SNAPSHOT.jar"
	err := os.Remove(dir)
	fmt.Println(dir)
	if err!=nil {
		fmt.Println(err)

	}
}