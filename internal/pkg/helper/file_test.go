package helper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	// Get this file directory path.
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	testPath := filepath.Join(dir, "/tmp")

	t.Run("ok", func(t *testing.T) {
		ioutil.WriteFile(testPath, []byte(""), 0644)
		exists := FileExists(testPath)
		assert.True(t, exists)
	})

	t.Run("ng", func(t *testing.T) {
		os.Remove(testPath)
		exists := FileExists(testPath)
		assert.False(t, exists)
	})
}
