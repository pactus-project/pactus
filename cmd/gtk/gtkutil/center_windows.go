//go:build gtk && windows

package gtkutil

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procSetWindowPos             = user32.NewProc("SetWindowPos")
	procGetSystemMetrics         = user32.NewProc("GetSystemMetrics")
)

type winRect struct {
	left, top, right, bottom int32
}

// CenterActiveWindow centers this process's foreground window on the primary
// monitor. GTK4 provides no way to move a window, so parentless/undecorated
// windows such as the splash are centered natively on Windows.
func CenterActiveWindow() {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return
	}

	var pid uint32
	procGetWindowThreadProcessID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid != windows.GetCurrentProcessId() {
		return
	}

	var rect winRect
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	width := rect.right - rect.left
	height := rect.bottom - rect.top

	screenW, _, _ := procGetSystemMetrics.Call(0) // SM_CXSCREEN
	screenH, _, _ := procGetSystemMetrics.Call(1) // SM_CYSCREEN

	x := (int32(screenW) - width) / 2
	y := (int32(screenH) - height) / 2

	const (
		swpNoSize     = 0x0001
		swpNoZOrder   = 0x0004
		swpNoActivate = 0x0010
	)
	procSetWindowPos.Call(hwnd, 0, uintptr(x), uintptr(y), 0, 0,
		swpNoSize|swpNoZOrder|swpNoActivate)
}
