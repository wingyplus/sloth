package sloth

import (
	"io/ioutil"
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

func TestCreateFileAtFirstWrite(t *testing.T) {
	makeTempDir("TestWrite")
	defer os.RemoveAll("TestWrite")

	logger := &Logger{
		Filename: "TestWrite/test-write.log",
	}

	logger.Write([]byte("Hello world"))

	if len(ls("TestWrite")) != 1 {
		t.Error("Expect has a file in folder TestWrite")
	}
}

func ls(dir string) []os.FileInfo {
	if infos, err := ioutil.ReadDir(dir); err == nil {
		return infos
	}
	return nil
}

func makeTempDir(dir string) {
	os.MkdirAll(dir, 0744)
}
