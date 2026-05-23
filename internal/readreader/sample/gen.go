package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// TODO: This is stupid remove this later
func main() {
	fmt.Println("Genrating the file")
	path := filepath.Join("./internal/readreader/sample", "file")
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Genration failed")
		return
	}
	defer file.Close()

	var size int64 = 1024
	err = file.Truncate(size)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Genration failed")
		return
	}
	fmt.Println("----------------Done----------------")
}
