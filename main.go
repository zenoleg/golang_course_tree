package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(writer io.Writer, path string, printFiles bool) error {
	rootFile, err := os.Open(path)

	if err != nil {
		return err
	}

	files, _ := rootFile.Readdir(-1)

	files = sortFiles(files)

	for iterator, file := range files {

		if !file.IsDir() && printFiles {
			_, _ = fmt.Fprintf(writer, "├───%s (%s)\n", file.Name(), processSize(file.Size()))
		} else if file.IsDir() {
			if iterator != len(files)-1 {
				_, _ = fmt.Fprintf(writer, "├───%s\n│\t", file.Name())
			} else {
				_, _ = fmt.Fprintf(writer, "└───%s\n", file.Name())
			}
			_ = dirTree(writer, path+"/"+file.Name(), printFiles)
		}
	}

	return nil
}

func processSize(size int64) string {
	if size == 0 {
		return "empty"
	}

	return strconv.FormatInt(size, 10) + "b"
}

func sortFiles(files []os.FileInfo) []os.FileInfo {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	return files
}
