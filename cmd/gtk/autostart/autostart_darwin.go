package autostart

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// launchAgentLabel matches the bundle identifier of Pactus GUI.
const launchAgentLabel = "com.github.pactus-project.pactus.pactus-gui"

// entryPath returns the launch agent of Pactus GUI.
func entryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Library", "LaunchAgents", launchAgentLabel+".plist"), nil
}

func executable() (string, error) {
	return os.Executable()
}

func enable(exe string, args []string) error {
	path, err := entryPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	return os.WriteFile(path, launchAgent(programArguments(exe, args)), 0o600)
}

func disable() error {
	path, err := entryPath()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}

func isEnabled() bool {
	path, err := entryPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(path)

	return err == nil
}

// programArguments returns the command of the launch agent. Inside an app
// bundle, the bundle is opened instead of the binary, so the launcher script
// can set up the GTK environment.
func programArguments(exe string, args []string) []string {
	const bundleMarker = ".app/Contents/MacOS/"

	idx := strings.Index(exe, bundleMarker)
	if idx < 0 {
		return append([]string{exe}, args...)
	}

	bundle := exe[:idx+len(".app")]
	command := []string{"/usr/bin/open", "-a", bundle}
	if len(args) > 0 {
		command = append(command, "--args")
		command = append(command, args...)
	}

	return command
}

func launchAgent(command []string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" ` +
		`"http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n")
	buf.WriteString("<plist version=\"1.0\">\n<dict>\n")
	buf.WriteString("\t<key>Label</key>\n\t<string>" + launchAgentLabel + "</string>\n")
	buf.WriteString("\t<key>ProgramArguments</key>\n\t<array>\n")
	for _, arg := range command {
		buf.WriteString("\t\t<string>")
		_ = xml.EscapeText(buf, []byte(arg))
		buf.WriteString("</string>\n")
	}
	buf.WriteString("\t</array>\n")
	buf.WriteString("\t<key>RunAtLoad</key>\n\t<true/>\n")
	buf.WriteString("</dict>\n</plist>\n")

	return buf.Bytes()
}
