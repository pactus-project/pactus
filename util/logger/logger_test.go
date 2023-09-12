package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type Foo struct{}

func (f Foo) String() string {
	return "foo"
}

func TestNilObjLogger(t *testing.T) {
	l := NewSubLogger("test", nil)
	var buf bytes.Buffer
	l.logger = l.logger.Output(&buf)

	l.Info("hello", "error", fmt.Errorf("error"))
	assert.Contains(t, buf.String(), "hello")
	assert.Contains(t, buf.String(), "error")
}

func TestObjLogger(t *testing.T) {
	globalInst = nil
	c := DefaultConfig()
	c.Colorful = false
	InitGlobalLogger(c)

	l := NewSubLogger("test", Foo{})
	var buf bytes.Buffer
	l.logger = l.logger.Output(&buf)

	l.Trace("a")
	l.Debug("b")
	l.Info("c")
	l.Warn("d")
	l.Error("e")

	out := buf.String()

	assert.Contains(t, out, "foo")
	assert.NotContains(t, out, "trace")
	assert.NotContains(t, out, "debug")
	assert.Contains(t, out, "info")
	assert.Contains(t, out, "warn")
	assert.Contains(t, out, "error")
}

func TestLogger(t *testing.T) {
	globalInst = nil
	c := DefaultConfig()
	c.Colorful = true
	InitGlobalLogger(c)

	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	Trace("a")
	Info("b", nil)
	Info("b", "a", nil)
	Info("c", "b", []byte{1, 2, 3})
	Warn("d", "x")
	Error("e", "y", Foo{})

	out := buf.String()

	fmt.Println(out)
	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "010203")
	assert.Contains(t, out, "!INVALID-KEY!")
	assert.Contains(t, out, "!MISSING-VALUE!")
	assert.Contains(t, out, "null")
	assert.NotContains(t, out, "trace")
	assert.NotContains(t, out, "debug")
	assert.Contains(t, out, "info")
	assert.Contains(t, out, "warn")
	assert.Contains(t, out, "error")
}

func TestRotating(t *testing.T) {
	tempDir := util.TempDirPath()
	fmt.Println(tempDir)
	MaxLogSize = 1
	LogFilename = filepath.Join(tempDir, "pactus.log")
	c := DefaultConfig()
	InitGlobalLogger(c)
	logger := NewSubLogger("test", nil)

	for i := 0; i < 1000; i++ {
		logger.Info(strings.Repeat("l", 1024))
	}

	assert.True(t, hasGzFile(tempDir), "log didn't rotate")
}

func hasGzFile(dir string) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".gz") {
			return true
		}
	}

	return false
}
