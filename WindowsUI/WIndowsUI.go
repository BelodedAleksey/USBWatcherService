package WindowsUI

import (
	"syscall"
)

const (
	IDABORT    = 3
	IDCANCEL   = 2
	IDCONTINUE = 11
	IDIGNORE   = 5
	IDNO       = 7
	IDOK       = 1
	IDRETRY    = 4
	IDTRYAGAIN = 10
	IDYES      = 6
)

var isInitialized bool = false
var kernel32 *syscall.DLL
var user32 *syscall.DLL
var wtsapi32 *syscall.DLL

func Init() {
	if isInitialized {
		return
	}
	isInitialized = true

	var err error
	kernel32, err = syscall.LoadDLL("kernel32.dll")
	if err != nil {
		panic(err)
	}
	user32, err = syscall.LoadDLL("user32.dll")
	if err != nil {
		panic(err)
	}
	wtsapi32, err = syscall.LoadDLL("wtsapi32.dll")
	if err != nil {
		panic(err)
	}
}
