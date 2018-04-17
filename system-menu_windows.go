package main

import (
	"syscall"

	"github.com/lxn/win"
)

var (
	user32Dll       = win.MustLoadLibrary("user32.dll")
	getSystemMenuFp = win.MustGetProcAddress(user32Dll, "GetSystemMenu")
)

// @see https://msdn.microsoft.com/en-us/library/ms647985
func GetSystemMenu(hWnd win.HWND, bRevert win.BOOL) win.HMENU {
	ret, _, _ := syscall.Syscall(getSystemMenuFp, 2,
		uintptr(hWnd),
		uintptr(bRevert),
		0)
	return win.HMENU(ret)
}



