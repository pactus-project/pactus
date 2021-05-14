package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
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
	if err := ioutil.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write to %s: %v", filename, err)
	}
	return nil
}

func Mkdir(dir string) error {
	// create the directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create directory %s", dir)
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

func TempDirPath() string {
	p, err := ioutil.TempDir("", "zarb*")
	if err != nil {
		panic(err)
	}
	return p
}

func TempFilePath() string {
	return path.Join(TempDirPath(), "file")
}

func IsDirEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	return err == io.EOF
}

func IsDirNotExistsOrEmpty(name string) bool {
	if !PathExists(name) {
		return true
	}

	return IsDirEmpty(name)
}

func IsValidDirPath(fp string) bool {
	fi, err := os.Stat(fp)
	if err == nil {
		if fi.IsDir() {
			if err := ioutil.WriteFile(fp+"/test", []byte{}, 0600); err != nil {
				return false
			}
			os.Remove(fp + "/test")
			return true
		}
		return false
	}

	if err := Mkdir(fp); err != nil {
		return false
	}
	os.Remove(fp)
	return true
}
