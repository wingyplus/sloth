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
	ext := filepath.Ext(logger.Filename)
	name := filepath.Base(logger.Filename)
	prefix := name[:len(name)-len(ext)]

	os.Create(fmt.Sprintf("%s_%s%s", prefix, timeNow().Format("20060102_1504"), ext))

	return nil
}
