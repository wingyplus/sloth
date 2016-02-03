package sloth

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// for mock in this package
var timeNow = time.Now

// Logger writes specific to Filename
type Logger struct {
	Filename string
}

func (logger *Logger) rotate() error {
	dir := filepath.Dir(logger.Filename)
	ext := filepath.Ext(logger.Filename)
	name := filepath.Base(logger.Filename)
	prefix := name[:len(name)-len(ext)]

	if !dirExist(dir) {
		os.MkdirAll(dir, 0744)
	}
	f, err := os.OpenFile(filepath.Join(dir, fmt.Sprintf("%s_%s%s", prefix, timeNow().Format("20060102_1504"), ext)), os.O_CREATE, 0644)
	defer f.Close()

	return err
}

func dirExist(dir string) bool {
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		return false
	}
	f.Close()
	return true
}
