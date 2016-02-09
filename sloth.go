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
var _ io.WriteCloser = (*File)(nil)

// Logger writes specific to Filename
type File struct {
	Filename string
	Every    time.Duration

	file      *os.File
	createdAt time.Time
}

func (f *File) Write(b []byte) (n int, err error) {
	if f.file == nil {
		f.file, err = openNew(f.Filename, false)
		f.createdAt = timeNow()
	}

	if now := timeNow(); f.Every > 0 && now.Sub(f.createdAt) >= f.Every {
		backup(f)
	}
	return f.file.Write(b)
}

func (f *File) Close() error {
	return f.file.Close()
}

func backup(f *File) {
	name := f.file.Name()
	f.Close()

	backupfile, err := openNew(name, true)
	if err != nil {
		return
	}
	defer backupfile.Close()

	f.file, _ = openNew(name, false)
	f.createdAt = timeNow()
	if _, err := io.Copy(backupfile, f.file); err != nil {
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
