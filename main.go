package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

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
	dirLastInfo := make(map[int]bool, 0)
	return tree(out, path, printFiles, 0, dirLastInfo)
}

func tree(out io.Writer, path string, printFiles bool, depth int, dirLastInfo map[int]bool) error {
	files, err := readFiles(path, printFiles)

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	filesCount := len(files) - 1
	var prefix string
	var walls string

	for position, file := range files {
		if position == filesCount {
			prefix = TAIL
			dirLastInfo[depth] = true
		} else {
			prefix = ARROW
			dirLastInfo[depth] = false
		}

		if depth > 0 {
			for i := 0; i < depth; i++ {
				if dirLastInfo[i] {
					fmt.Fprint(out, "\t")
				} else {
					fmt.Fprint(out, WALL+"\t")
				}
			}

			fmt.Fprint(out, walls)
		}

		if file.IsDir() {
			fmt.Fprint(out, formatDir(file, prefix))

			fullPath := path + "/" + file.Name()

			err := tree(out, fullPath, printFiles, depth+1, dirLastInfo)

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

func formatDir(file os.FileInfo, prefix string) string {
	return fmt.Sprintf("%s%s\n", prefix, file.Name())
}
