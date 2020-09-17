package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"proteinreminder/internal/pkg/testutil"
	"runtime"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Get this file directory path.
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	testPath := filepath.Join(dir, "/tmp")
	ioutil.WriteFile(testPath, []byte(""), 0644)

	exists := FileExists(testPath)
	if exists != true {
		t.Errorf(testutil.MakeTestMessageWithGotWant(exists, true))
	}

	os.Remove(testPath)

	exists = FileExists(testPath)
	if exists != false {
		t.Errorf(testutil.MakeTestMessageWithGotWant(exists, false))
	}
}
