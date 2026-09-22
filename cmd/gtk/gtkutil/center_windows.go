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

// callProc invokes a Win32 entry point and returns its result. The Call error
// is deliberately ignored: these entry points report failure through the
// return value, not through GetLastError.
func callProc(proc *windows.LazyProc, args ...uintptr) uintptr {
	ret, _, _ := proc.Call(args...)

	return ret
}

// CenterActiveWindow centers this process's foreground window on the primary
// monitor. GTK4 provides no way to move a window, so parentless/undecorated
// windows such as the splash are centered natively on Windows.
func CenterActiveWindow() {
	hwnd := callProc(procGetForegroundWindow)
	if hwnd == 0 {
		return
	}

	var pid uint32
	//nolint:gosec // the Win32 out-parameter requires an unsafe pointer
	callProc(procGetWindowThreadProcessID, hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid != windows.GetCurrentProcessId() {
		return
	}

	var rect winRect
	//nolint:gosec // the Win32 out-parameter requires an unsafe pointer
	callProc(procGetWindowRect, hwnd, uintptr(unsafe.Pointer(&rect)))
	width := rect.right - rect.left
	height := rect.bottom - rect.top

	screenW := callProc(procGetSystemMetrics, 0) // SM_CXSCREEN
	screenH := callProc(procGetSystemMetrics, 1) // SM_CYSCREEN

	posX := (int32(screenW) - width) / 2
	posY := (int32(screenH) - height) / 2

	const (
		swpNoSize     = 0x0001
		swpNoZOrder   = 0x0004
		swpNoActivate = 0x0010
	)
	callProc(procSetWindowPos, hwnd, 0, uintptr(posX), uintptr(posY), 0, 0,
		swpNoSize|swpNoZOrder|swpNoActivate)
}
