package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	StickSeparator  = "│\t"
	CommonSeparator = "\t"
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
	return tree(writer, path, printFiles, "")
}

func tree(writer io.Writer, path string, printFiles bool, lastStr string) error {
	files, err := readFiles(path, printFiles)

	if err != nil {
		return err
	}

	dirDepth := strings.Count(path, "/")

	var firstSeparator, lastString string

	stickCount := calculateStickCount(lastStr)

	if stickCount > 0 {
		firstSeparator = StickSeparator
	}

	for iterator, file := range files {
		lastString = ""

		if dirDepth >= 1 {
			prefix := strings.Repeat(firstSeparator, stickCount) + strings.Repeat(CommonSeparator, dirDepth-stickCount)

			_, _ = fmt.Fprintf(writer, prefix)
			lastString += prefix
		}

		var branch string

		if file.IsDir() {
			if iterator == len(files)-1 {
				branch = fmt.Sprintf("└───%s\n", file.Name())
			} else {
				branch = fmt.Sprintf("├───%s\n", file.Name())
			}

			_, _ = fmt.Fprintf(writer, branch)
			lastString += branch

			newPath := path + "/" + file.Name()
			if dirCount, _ := fileCount(newPath); dirCount > 0 {
				_ = tree(writer, newPath, printFiles, lastString)
			}
		} else {
			if iterator == len(files)-1 {
				branch = fmt.Sprintf("└───%s (%s)\n", file.Name(), processSize(file.Size()))
			} else {
				branch = fmt.Sprintf("├───%s (%s)\n", file.Name(), processSize(file.Size()))
			}

			_, _ = fmt.Fprintf(writer, branch)
			lastString += branch
		}
	}

	return nil
}

func calculateStickCount(lastStr string) int {
	stickCount := strings.Count(lastStr, "│")
	arrowCount := strings.Count(lastStr, "├───")

	return stickCount + arrowCount
}

func readFiles(path string, printFiles bool) ([]os.FileInfo, error) {
	rootFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	files, _ := rootFile.Readdir(-1)
	files = sortFiles(files)

	var fileList []os.FileInfo

	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, file)
		} else if printFiles {
			fileList = append(fileList, file)
		}
	}

	return fileList, nil
}

func fileCount(path string) (int, error) {
	rootFile, err := os.Open(path)

	if err != nil {
		return 0, err
	}

	files, _ := rootFile.Readdir(-1)

	return len(files), nil
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
