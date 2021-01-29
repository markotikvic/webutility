package webutility

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// ReadFileContent ...
func ReadFileContent(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	if _, err = io.Copy(buf, f); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ReadFileLines ...
func ReadFileLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s strings.Builder

	if _, err = io.Copy(&s, f); err != nil {
		return nil, err
	}

	lines := strings.Split(s.String(), "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r\n")
	}

	return lines, nil
}

// LinesToFile ...
func LinesToFile(path string, lines []string) error {
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}

	return ioutil.WriteFile(path, []byte(content), 0644) // drw-r--r--
}

// InsertLine ...
func InsertLine(lines *[]string, pos int64, l string) {
	tail := append([]string{l}, (*lines)[pos:]...)

	*lines = append((*lines)[:pos], tail...)
}

func WriteFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0644) // drw-r--r--
}

func ListDir(path string) (fnames []string, err error) {
	finfo, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range finfo {
		fnames = append(fnames, f.Name())
	}

	return fnames, nil
}

func WorkingDir() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("couldn't get working directory: %s\n", err.Error())
	}
	return path
}

func FileExtension(path string) string {
	parts := strings.Split(path, ".") // because name can contain dots
	if len(parts) < 2 {
		return ""
	}
	return "." + parts[len(parts)-1]
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

// DirectoryFromPath ...
func DirectoryFromPath(path string) (dir string) {
	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return ""
	}

	dir = parts[0]
	for _, p := range parts[1 : len(parts)-1] {
		dir += "/" + p
	}

	return dir
}

// FileExists ...
func FileExists(path string) bool {
	temp, err := os.Open(path)
	defer temp.Close()

	if err != nil {
		return false
	}

	return true
}
