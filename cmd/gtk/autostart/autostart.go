// Package autostart manages the entry that starts Pactus GUI when the user
// logs in to the operating system.
package autostart

import (
	"errors"
	"flag"
	"fmt"
	"slices"
)

// ErrUnsupported is returned when autostart is not available on this platform.
var ErrUnsupported = errors.New("autostart is not supported on this platform")

// Entry is the login entry that starts Pactus GUI with the given arguments.
type Entry struct {
	args []string
}

// New creates an autostart entry that starts Pactus GUI with the given arguments.
func New(args []string) *Entry {
	return &Entry{args: args}
}

// IsEnabled reports whether the autostart entry exists.
func (*Entry) IsEnabled() bool {
	return isEnabled()
}

// SetEnabled creates or removes the autostart entry.
func (e *Entry) SetEnabled(enabled bool) error {
	if !enabled {
		return disable()
	}

	exe, err := executable()
	if err != nil {
		return err
	}

	return enable(exe, e.args)
}

// ArgsFromFlags returns the flags that were set on the command line, in the
// "-name=value" form, skipping the excluded ones. Secrets such as the wallet
// password must be excluded, because the autostart entry is stored in plain text.
func ArgsFromFlags(flagSet *flag.FlagSet, exclude ...string) []string {
	args := []string{}
	flagSet.Visit(func(f *flag.Flag) {
		if slices.Contains(exclude, f.Name) {
			return
		}
		args = append(args, fmt.Sprintf("-%s=%s", f.Name, f.Value.String()))
	})

	return args
}
