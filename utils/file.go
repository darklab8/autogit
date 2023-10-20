package utils

import (
	"autogit/settings/logus"
	"autogit/settings/types"
	"bufio"
	"fmt"
	"os"
)

type myfile struct {
	Filepath types.FilePath
	file     *os.File
	lines    []string
}

func NewFile(Filepath types.FilePath) *myfile {
	return &myfile{
		Filepath: Filepath,
	}
}

func (f *myfile) OpenToReadF() *myfile {
	file, err := os.Open(string(f.Filepath))
	f.file = file

	logus.CheckFatal(err, "failed to open file", logus.FilePath(f.Filepath))
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
	file, err := os.Create(string(f.Filepath))
	logus.CheckFatal(err, "failed to create file ", logus.FilePath(f.Filepath))
	f.file = file
	return f
}
func (f *myfile) WritelnF(msg string) {
	_, err := f.file.WriteString(fmt.Sprintf("%v\n", msg))

	logus.CheckFatal(err, "failed to write string to file")
}

func FileExists(filename types.FilePath) bool {
	info, err := os.Stat(string(filename))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
