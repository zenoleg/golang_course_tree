package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

/**
├───project
│	├───file.txt (19b)
│	└───gopher.png (70372b)
├───static
│	├───a_lorem
│	│	├───dolor.txt (empty)
│	│	├───gopher.png (70372b)
│	│	└───ipsum
│	│		└───gopher.png (70372b)
│	├───css
│	│	└───body.css (28b)
│	├───empty.txt (empty)
│	├───html
│	│	└───index.html (57b)
│	├───js
│	│	└───site.js (10b)
│	└───z_lorem
│		├───dolor.txt (empty)
│		├───gopher.png (70372b)
│		└───ipsum
│			└───gopher.png (70372b)
├───zline
│	├───empty.txt (empty)
│	└───lorem
│		├───dolor.txt (empty)
│		├───gopher.png (70372b)
│		└───ipsum
│			└───gopher.png (70372b)
└───zzfile.txt (empty)
*/

const (
	ARROW = "├───"
	WALL  = "│"
	TAIL  = "└───"
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

func dirTree(out io.Writer, path string, printFiles bool) error {
	files, err := readFiles(path, printFiles)

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	filesCount := len(files) - 1
	var prefix string

	for position, file := range files {
		if position == filesCount {
			prefix = TAIL
		} else {
			prefix = ARROW
		}

		if file.IsDir() {
			fmt.Fprint(out, formatFile(file, prefix))

			fullPath := path + "/" + file.Name()

			err := dirTree(out, fullPath, printFiles)

			if err != nil {
				return err
			}
		} else {
			fmt.Fprint(out, formatFile(file, prefix))
		}
	}

	return nil
}

func readFiles(path string, printFiles bool) ([]os.FileInfo, error) {
	root, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	files, err := root.Readdir(-1)
	if err != nil {
		return nil, err
	}

	files = sortFiles(files)

	fileList := make([]os.FileInfo, 0, 20)

	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, file)
		} else if printFiles {
			fileList = append(fileList, file)
		}
	}

	return fileList, nil
}

func sortFiles(files []os.FileInfo) []os.FileInfo {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	return files
}

func formatFile(file os.FileInfo, prefix string) string {
	var size string

	if file.Size() == 0 {
		size = "empty"
	} else {
		size = strconv.FormatInt(file.Size(), 10) + "b"
	}

	return fmt.Sprintf("%s%s (%s)\n", prefix, file.Name(), size)
}
