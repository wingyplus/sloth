package sloth

import (
	"os"
	"testing"
	"time"
)

func init() {
	timeNow = func() time.Time {
		return time.Date(2016, time.February, 3, 15, 40, 0, 0, time.UTC)
	}
}

func TestRotate(t *testing.T) {
	logger := &Logger{
		Filename: "./test.log",
	}

	logger.rotate()

	f, err := os.Open("./test_20160203_1540.log")
	if os.IsNotExist(err) {
		t.Error("Expect file ./test_20160203_1540.log is exist.")
	}
	f.Close()
}
