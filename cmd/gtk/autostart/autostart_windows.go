package autostart

import (
	"errors"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

const runValueName = "PactusGUI"

// runKeyPath is the per-user registry key of programs started at logon.
// It is a variable so tests can use a temporary key.
var runKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`

func executable() (string, error) {
	return os.Executable()
}

func enable(exe string, args []string) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func() { _ = key.Close() }()

	return key.SetStringValue(runValueName, commandLine(exe, args))
}

func disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer func() { _ = key.Close() }()

	err = key.DeleteValue(runValueName)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}

	return err
}

func isEnabled() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer func() { _ = key.Close() }()

	_, _, err = key.GetStringValue(runValueName)

	return err == nil
}

func commandLine(exe string, args []string) string {
	parts := make([]string, 0, len(args)+1)
	for _, arg := range append([]string{exe}, args...) {
		parts = append(parts, syscall.EscapeArg(arg))
	}

	return strings.Join(parts, " ")
}
