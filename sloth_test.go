package sloth

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var (
	oldTime = time.Date(2016, time.February, 3, 15, 40, 0, 0, time.UTC)
	newTime = time.Date(2016, time.February, 3, 15, 45, 0, 0, time.UTC)

	currentTime = &oldTime
)

func init() {
	timeNow = func() time.Time {
		return *currentTime
	}
}

var logRotateTestCases = []struct {
	logger   *Logger
	filename string
}{
	{
		&Logger{Filename: "./TestRotate/test.log"}, "./TestRotate/test_20160203_1540.log",
	},
	{
		&Logger{Filename: "./TestRotate/InFolder/test.log"}, "./TestRotate/InFolder/test_20160203_1540.log",
	},
}

func TestRotate(t *testing.T) {
	makeTempDir("TestRotate")
	defer os.RemoveAll("TestRotate")

	for _, testcase := range logRotateTestCases {
		testcase.logger.rotate()

		if !exist(testcase.filename) {
			t.Errorf("Expect file %s is exist.", testcase.filename)
		}
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

func TestAutoRotate(t *testing.T) {
	makeTempDir("TestAutoRotate")
	defer os.RemoveAll("TestAutoRotate")

	logger := &Logger{
		Filename: "TestAutoRotate/test-auto-rotate.log",
		Every:    1 * time.Millisecond,
	}

	logger.Write([]byte("Hello world"))
	updateTime()
	logger.Write([]byte("Hello world"))

	if total := len(ls("TestAutoRotate")); total != 2 {
		t.Error("Expect have 2 file in folder TestWrite but got", total)
	}

	resetTime()
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

func updateTime() {
	*currentTime = newTime
}

func resetTime() {
	*currentTime = oldTime
}
