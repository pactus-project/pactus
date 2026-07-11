//go:build gtk && !windows

package main

// forceEnglishUILanguage is a no-op outside Windows, where the LANGUAGE
// environment variable already controls GTK's message catalogs.
func forceEnglishUILanguage() {}
