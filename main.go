package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
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
	return myDirTree(writer, path, printFiles, false)
}

func myDirTree(writer io.Writer, path string, printFiles bool, isPrevLast bool) error {
	files, err := readOnlyDirectories(path)
	if err != nil {
		return err
	}

	dirDepth := strings.Count(path, "/")

	var separator string

	if isPrevLast {
		separator = "\t"
	} else {
		separator = "│\t"
	}

	for iterator, file := range files {

		if file.Name() == ".git" || file.Name() == ".idea" {
			continue
		}

		if file.IsDir() {
			if dirDepth >= 1 {
				_, _ = fmt.Fprintf(writer, strings.Repeat(separator, dirDepth))
			}

			if iterator == len(files)-1 {
				_, _ = fmt.Fprintf(writer, "└───%s\n", file.Name())
			} else {
				_, _ = fmt.Fprintf(writer, "├───%s\n", file.Name())
			}

			newPath := path + "/" + file.Name()
			if dirCount, _ := dirCount(newPath); dirCount > 0 {
				_ = myDirTree(writer, newPath, printFiles, iterator == len(files)-1)
			}
		}
	}

	return nil
}

func readOnlyDirectories(path string) ([]os.FileInfo, error) {
	rootFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	files, _ := rootFile.Readdir(-1)
	files = sortFiles(files)

	var dirFiles []os.FileInfo

	for _, file := range files {
		if file.IsDir() {
			dirFiles = append(dirFiles, file)
		}
	}

	return dirFiles, nil
}

func dirCount(path string) (int, error) {
	counter := 0

	rootFile, err := os.Open(path)

	if err != nil {
		return 0, err
	}

	files, _ := rootFile.Readdir(-1)

	for _, file := range files {
		if file.IsDir() {
			counter++
		}
	}

	return counter, nil
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
