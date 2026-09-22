//go:build gtk && !windows

package gtkutil

// CenterActiveWindow is a no-op on platforms other than Windows, where the
// window manager already positions new toplevels reasonably.
func CenterActiveWindow() {}
