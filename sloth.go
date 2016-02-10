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

// File wrap os.File to manage rotate logic
type File struct {
	Filename string
	Every    time.Duration

	file      *os.File
	createdAt time.Time
}

// Write data to the file
func (f *File) Write(b []byte) (n int, err error) {
	if f.file == nil {
		if err = f.openNew(); err != nil {
			return
		}
	}

	if now := timeNow(); f.Every > 0 && now.Sub(f.createdAt) >= f.Every {
		if err = f.backup(); err != nil {
			return
		}

		if err = f.openNew(); err != nil {
			return
		}
	}
	return f.file.Write(b)
}

// Close the file
func (f *File) Close() error {
	return f.file.Close()
}

// Name return file name
func (f *File) Name() string { return f.file.Name() }

func (f *File) backup() error {
	name := f.Name()
	f.Close()

	return os.Rename(name, filepath.Join(filepath.Dir(name), backupName(name)))
}

func (f *File) openNew() (err error) {
	f.file, err = os.OpenFile(f.Filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	f.createdAt = timeNow()
	return
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func backupName(filename string) string {
	ext := filepath.Ext(filename)
	name := filepath.Base(filename)
	prefix := name[:len(name)-len(ext)]

	return fmt.Sprintf("%s_%s%s", prefix, timeNow().Format("20060102_1504"), ext)
}
