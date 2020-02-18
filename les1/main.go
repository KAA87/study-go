package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	PART1 = "├───"
	PART2 = "└───"
	PART3 = "	"
	PART4 = "│	"
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
		panic(err)
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("not find path")
	}

	typeFile, err := checkTypeFile(absPath)
	if err != nil {
		return fmt.Errorf("error check type file")
	}

	if typeFile == "file" && printFiles == false {
		return nil
	}

	level := len(strings.Split(absPath, string(os.PathSeparator)))
	err = formattedDirTree(out, absPath, printFiles, level, 0, false, []string{})
	if err != nil {
		return fmt.Errorf("")
	}

	return nil
}

func formattedDirTree(out io.Writer, path string, printFiles bool, level int, baseLevel int, flLast bool, padding []string) error {

	typeFile, err := checkTypeFile(path)
	if err != nil {
		return fmt.Errorf("error check type file")
	}
	if typeFile == "file" && printFiles == false {
		return nil
	}

	if baseLevel != 0 {

		itemLevel := PART1
		if flLast == true {
			itemLevel = PART2
		}

		padding = append(padding, itemLevel)

		padding = formattedPadding(padding)

		pathArr := strings.Split(path, string(os.PathSeparator))

		sizeF := ""
		if printFiles == true && typeFile == "file" {
			sizeF = " (empty)"
			if sizeFile := getSizeFile(path); sizeFile > 0 {
				sizeF = " (" + strconv.Itoa(sizeFile) + "b)"
			}
		}

		strOut := strings.Join(padding, "") + "" + pathArr[level+baseLevel-1] + sizeF

		fmt.Fprintln(out, strOut)
	}

	if typeFile == "dir" {

		baseLevel++

		files, err := filepath.Glob(path + "/*")
		if err != nil {
			fmt.Print(err)
			return fmt.Errorf("err not get list file")
		}

		files = prepareFiles(files, printFiles)

		lenFiles := len(files)
		if lenFiles > 0 {
			for iFP, filePath := range files {

				flLast = false
				if iFP == lenFiles-1 {
					flLast = true
				}

				err := formattedDirTree(out, filePath, printFiles, level, baseLevel, flLast, padding)
				if err != nil {
					fmt.Print(err)
					return fmt.Errorf("err dir tree")
				}
			}
		}

	}

	return nil
}

// проверка на директорию или файл
func checkTypeFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("err open file")
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Print(err)
		return "", fmt.Errorf("err get info file")
	}

	if fi.IsDir() {
		return "dir", nil
	}
	return "file", nil
}

// получить размер файла
func getSizeFile(path string) int {
	fi, _ := os.Stat(path)

	size := int(fi.Size())

	return size
}

// подготовка файлов, сортировка и удаление всех файлов если не задан параметр -f
func prepareFiles(files []string, printFiles bool) []string {

	if printFiles == false {
		for k, i := range files {
			typeFile, _ := checkTypeFile(i)

			if typeFile == "file" {
				files[k] = files[len(files)-1]
				files[len(files)-1] = ""
				files = files[:len(files)-1]
			}
		}
	}

	sort.Strings(files)

	return files
}

// форматируем отступы опираясь на текущий
func formattedPadding(padding []string) []string {
	lenPadding := len(padding)
	if lenPadding >= 2 {
		switch padding[lenPadding-2] {
		case PART2:
			padding[lenPadding-2] = PART3
		case PART1:
			padding[lenPadding-2] = PART4
		}
	}

	return padding
}
