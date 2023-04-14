package utils

import (
	"bufio"
	"fmt"
	"os"
)

type myfile struct {
	Filepath string
	file     *os.File
	lines    []string
}

func NewFile(Filepath string) *myfile {
	return &myfile{
		Filepath: Filepath,
	}
}

func (f *myfile) OpenToReadF() *myfile {
	file, err := os.Open(f.Filepath)
	f.file = file

	CheckFatal(err, "failed to open ", f.Filepath)
	return f
}

func (f *myfile) Close() {
	f.file.Close()
}

func (f *myfile) ReadLines() []string {

	scanner := bufio.NewScanner(f.file)

	for scanner.Scan() {
		f.lines = append(f.lines, scanner.Text())
	}
	return f.lines
}

func (f *myfile) FileReadLines() []string {
	f.OpenToReadF()
	defer f.Close()
	return f.ReadLines()
}

func (f *myfile) CreateToWriteF() *myfile {
	file, err := os.Create(f.Filepath)
	f.file = file
	CheckFatal(err, "failed to open ", f.Filepath)

	return f
}
func (f *myfile) WritelnF(msg string) {
	_, err := f.file.WriteString(fmt.Sprintf("%v\n", msg))

	CheckFatal(err, "failed to write string to file")
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
