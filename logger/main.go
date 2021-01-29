package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const dateTimeFormat = "2006-01-02 15:04:05"

// Block ...
const (
	MaxLogSize5MB   int64 = 5 * 1024 * 1024
	MaxLogSize1MB   int64 = 1 * 1024 * 1024
	MaxLogSize500KB int64 = 500 * 1024
	MaxLogSize100KB int64 = 100 * 1024
	MaxLogSize512B  int64 = 512
)

// Logger ...
type Logger struct {
	mu         *sync.Mutex
	outputFile *os.File

	outputFileName string
	fullName       string
	maxFileSize    int64

	splitCount int

	directory string
}

// New ...
func New(name, dir string, maxFileSize int64) (logger *Logger, err error) {
	logger = &Logger{}

	logger.outputFileName = name
	logger.mu = &sync.Mutex{}
	logger.maxFileSize = maxFileSize
	logger.directory = dir

	err = os.Mkdir(dir, os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	date := strings.Replace(time.Now().Format(dateTimeFormat), ":", ".", -1)
	logger.outputFileName += " " + date
	logger.fullName = logger.outputFileName + ".txt"
	path := filepath.Join(dir, logger.fullName)
	if logger.outputFile, err = os.Create(path); err != nil {
		return nil, err
	}

	return logger, nil
}

// Log ...
func (l *Logger) Log(format string, v ...interface{}) {
	if l.outputFile == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, v...)
	s := time.Now().Format(dateTimeFormat) + ": " + msg + "\n"
	if l.shouldSplit(len(s)) {
		l.split()
	}
	l.outputFile.WriteString(s)
}

// Print ...
func (l *Logger) Print(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Printf("%s: %s\n", time.Now().Format(dateTimeFormat), msg)
}

// Trace ...
func (l *Logger) Trace(format string, v ...interface{}) string {
	if l.outputFile == nil {
		return ""
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	s := getTrace(format, v...) + "\n"

	if l.shouldSplit(len(s)) {
		l.split()
	}
	l.outputFile.WriteString(s)

	return s
}

// PrintTrace ...
func (l *Logger) PrintTrace(format string, v ...interface{}) {
	s := getTrace(format, v...)
	fmt.Printf("%s\n", s)
}

// PrintAndTrace ...
func (l *Logger) PrintAndTrace(format string, v ...interface{}) {
	t := l.Trace(format, v...)
	fmt.Println(t)
}

// Close ...
func (l *Logger) Close() error {
	if l.outputFile == nil {
		return nil
	}

	return l.outputFile.Close()
}

// GetOutDir ...
func (l *Logger) GetOutDir() string {
	return l.directory
}

func (l *Logger) split() error {
	if l.outputFile == nil {
		return nil
	}

	// close old file
	err := l.outputFile.Close()
	if err != nil {
		return err
	}

	// open new file
	l.splitCount++
	path := filepath.Join(l.directory, l.outputFileName+fmt.Sprintf("(%d)", l.splitCount)+".txt")
	l.outputFile, err = os.Create(path)

	return err
}

func (l *Logger) shouldSplit(nextEntrySize int) bool {
	stats, _ := l.outputFile.Stat()
	return int64(nextEntrySize) >= (l.maxFileSize - stats.Size())
}

func printJSON(in []byte) (out []byte, err error) {
	var buf bytes.Buffer
	err = json.Indent(&buf, in, "", "    ")
	return buf.Bytes(), err
}
