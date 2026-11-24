package autostart

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// entryPath returns the XDG autostart desktop entry of Pactus GUI.
//
// Inside a Flatpak, XDG_CONFIG_HOME points to the private config folder of the
// application, which the desktop session never reads. The entry is then written
// to the config folder of the user instead, which the sandbox can reach through
// its access to the home folder.
func entryPath() (string, error) {
	var dir string
	if os.Getenv("FLATPAK_ID") != "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".config")
	} else {
		var err error
		dir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(dir, "autostart", "pactus-gui.desktop"), nil
}

// executable returns the file to start. When running from an AppImage, the
// binary lives in a temporary mount, so the AppImage file is used instead.
// Inside a Flatpak, the binary path exists only in the sandbox, so the
// application is started with "flatpak" (see enable).
func executable() (string, error) {
	if os.Getenv("FLATPAK_ID") != "" {
		return "flatpak", nil
	}

	if appImage := os.Getenv("APPIMAGE"); appImage != "" {
		return appImage, nil
	}

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

	if id := os.Getenv("FLATPAK_ID"); id != "" {
		args = append([]string{"run", id}, args...)
	}

	return os.WriteFile(path, []byte(desktopEntry(exe, args)), 0o600)
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

func desktopEntry(exe string, args []string) string {
	fields := make([]string, 0, len(args)+1)
	for _, arg := range append([]string{exe}, args...) {
		fields = append(fields, quoteExecArg(arg))
	}

	return "[Desktop Entry]\n" +
		"Type=Application\n" +
		"Name=Pactus GUI\n" +
		"Exec=" + strings.Join(fields, " ") + "\n" +
		"Terminal=false\n" +
		"X-GNOME-Autostart-enabled=true\n"
}

// quoteExecArg quotes an argument for the Exec key of a desktop entry,
// following the Desktop Entry Specification: reserved characters are escaped
// inside the quotes, then backslashes are escaped again for the string value,
// and "%" is doubled so it is not read as a field code.
func quoteExecArg(arg string) string {
	quoted := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "`", "\\`", `$`, `\$`).Replace(arg)
	quoted = strings.ReplaceAll(quoted, `\`, `\\`)
	quoted = strings.ReplaceAll(quoted, `%`, `%%`)

	return `"` + quoted + `"`
}
