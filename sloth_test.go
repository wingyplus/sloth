package sloth

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func init() {
	resetTime()
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
		defer testcase.logger.Close()

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
	defer logger.Close()

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
	defer logger.Close()

	logger.Write([]byte("Hello world"))
	updateTime()
	logger.Write([]byte("Hello world"))

	if total := len(ls("TestAutoRotate")); total != 2 {
		t.Error("Expect have 2 file in folder TestWrite but got", total)
	}
	if s := cat("TestAutoRotate/test-auto-rotate_20160203_1545.log"); s != "Hello world" {
		t.Error("Expect `Hello World` in file content but got", s)
	}

	resetTime()
}

func TestWrite(t *testing.T) {
	makeTempDir("TestWrite")
	defer os.RemoveAll("TestWrite")

	logger := &Logger{
		Filename: "TestWrite/test-write.log",
		Every:    1 * time.Millisecond,
	}
	defer logger.Close()

	_, err := logger.Write([]byte("Hello world"))

	if err != nil {
		t.Error(err)
	}
	if s := cat("TestWrite/test-write.log"); s != "Hello world" {
		t.Error("Expect `Hello World` in file content but got", s)
	}
}

func cat(filename string) string {
	if b, err := ioutil.ReadFile(filename); err == nil {
		return string(b)
	}
	return ""
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
	timeNow = func() time.Time {
		return time.Date(2016, time.February, 3, 15, 45, 0, 0, time.UTC)
	}
}

func resetTime() {
	timeNow = func() time.Time {
		return time.Date(2016, time.February, 3, 15, 40, 0, 0, time.UTC)
	}
}
