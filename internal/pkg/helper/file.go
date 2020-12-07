// Package helper provides shared helper function.
package helper

import (
	"os"
	"path/filepath"
)

// GetAppDir returns the directory located the app binary. An error is returned if it cannot get the app executable path.
func GetAppDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

// FileExists check if filePath exists.
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
