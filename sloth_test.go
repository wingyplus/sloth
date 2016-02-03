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

var logRotateTestCases = []struct {
	logger   *Logger
	filename string
}{
	{
		&Logger{Filename: "./test.log"}, "./test_20160203_1540.log",
	},
	{
		&Logger{Filename: "./test/test.log"}, "./test/test_20160203_1540.log",
	},
}

func TestRotate(t *testing.T) {

	for _, testcase := range logRotateTestCases {
		testcase.logger.rotate()

		f, err := os.Open(testcase.filename)
		if os.IsNotExist(err) {
			t.Errorf("Expect file %s is exist.", testcase.filename)
		}
		f.Close()
	}
}
