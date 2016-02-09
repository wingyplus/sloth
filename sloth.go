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
var _ io.WriteCloser = (*Logger)(nil)

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
		backup(logger)
	}
	return logger.file.Write(b)
}

func (logger *Logger) Close() error {
	return logger.file.Close()
}

func backup(logger *Logger) {
	name := logger.file.Name()
	logger.Close()

	backupfile, err := openNew(name, true)
	if err != nil {
		return
	}
	defer backupfile.Close()

	logger.file, _ = openNew(name, false)
	if _, err := io.Copy(backupfile, logger.file); err != nil {
		panic(err)
	}
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func openNew(filename string, stamptime bool) (*os.File, error) {
	if !stamptime {
		return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	}

	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	name := filepath.Base(filename)
	prefix := name[:len(name)-len(ext)]

	if !exist(dir) {
		os.MkdirAll(dir, 0744)
	}

	return os.OpenFile(filepath.Join(dir, fmt.Sprintf("%s_%s%s", prefix, timeNow().Format("20060102_1504"), ext)), os.O_CREATE|os.O_WRONLY, 0644)
}
