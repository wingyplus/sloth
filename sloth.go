package sloth

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// for mock in this package
var timeNow = time.Now

// Make sure Logger always implements io.Writer
var _ io.Writer = (*Logger)(nil)

// Logger writes specific to Filename
type Logger struct {
	Filename string
	Every    time.Duration

	file      *os.File
	createdAt time.Time
}

func (logger *Logger) Write(b []byte) (n int, err error) {
	if logger.file == nil {
		logger.file, err = openNew(logger.Filename, false)
		logger.createdAt = timeNow()
	}

	if now := timeNow(); logger.Every > 0 && now.Sub(logger.createdAt) >= logger.Every {
		backup(logger.file)
	}
	return
}

func (logger *Logger) rotate() error {
	f, err := openNew(logger.Filename, true)
	defer f.Close()
	return err
}

func backup(logfile *os.File) {
	backupfile, err := openNew(logfile.Name(), true)
	if err != nil {
		return
	}
	defer backupfile.Close()
}

func dirExist(dir string) bool {
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		return false
	}
	f.Close()
	return true
}

func openNew(filename string, stamptime bool) (*os.File, error) {
	if !stamptime {
		return os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0644)
	}

	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	name := filepath.Base(filename)
	prefix := name[:len(name)-len(ext)]

	if !dirExist(dir) {
		os.MkdirAll(dir, 0744)
	}

	return os.OpenFile(filepath.Join(dir, fmt.Sprintf("%s_%s%s", prefix, timeNow().Format("20060102_1504"), ext)), os.O_CREATE, 0644)
}
