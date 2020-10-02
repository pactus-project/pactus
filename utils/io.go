package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func IsAbsPath(path string) bool {
	return filepath.IsAbs(path)
}

func MakeAbs(path string) string {
	if IsAbsPath(path) {
		return path
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Clean(filepath.Join(wd, path))
}

func ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func WriteFile(filename string, data []byte) error {
	// create directory
	if err := Mkdir(filepath.Dir(filename)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0777); err != nil {
		return fmt.Errorf("Failed to write to %s: %v", filename, err)
	}
	return nil
}

func Mkdir(dir string) error {
	// create the directory
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("Could not create directory %s", dir)
	}
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func TempPath() string {
	p, err := ioutil.TempDir("", "zarb*")
	if err != nil {
		panic(err)
	}
	return p
}
