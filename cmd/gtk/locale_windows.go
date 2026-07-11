//go:build gtk && windows

package main

import "golang.org/x/sys/windows"

// forceEnglishUILanguage sets the thread UI language to English (en-US), so
// GTK and libadwaita load their English catalogs instead of following the
// Windows display language. On Windows GTK reads this rather than the LANGUAGE
// environment variable.
func forceEnglishUILanguage() {
	const langEnUS = 0x0409

	proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("SetThreadUILanguage")
	_, _, _ = proc.Call(uintptr(langEnUS))
}
