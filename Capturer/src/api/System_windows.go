package api

import (
	"syscall"
	"unsafe"
)

var (
	mStdHandle                   uintptr
	mStdColor                    uintptr
	mUser32                      *syscall.DLL
	mKernel32                    *syscall.DLL
	mRemoveMenuProc              *syscall.Proc
	mGetStdHandleProc            *syscall.Proc
	mGetSystemMenuProc           *syscall.Proc
	mSetConsoleTitleWProc        *syscall.Proc
	mGetConsoleWindowProc        *syscall.Proc
	mSetConsoleTextAttributeProc *syscall.Proc
)

func init() {
	mUser32, _ = syscall.LoadDLL("User32.dll")
	mKernel32, _ = syscall.LoadDLL("Kernel32.dll")
	mGetStdHandleProc, _ = mKernel32.FindProc("GetStdHandle")
	mSetConsoleTitleWProc, _ = mKernel32.FindProc("SetConsoleTitleW")
	mGetConsoleWindowProc, _ = mKernel32.FindProc("GetConsoleWindow")
	mSetConsoleTextAttributeProc, _ = mKernel32.FindProc("SetConsoleTextAttribute")
	mStdHandle, _, _ = mGetStdHandleProc.Call(0xfffffff5)
}

// SetConsoleTitle 设置控制台标题
func SetConsoleTitle(title string) {
	if mSetConsoleTitleWProc != nil {
		mSetConsoleTitleWProc.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))))
	}
}

// SetConsoleTextColor 设置控制台文本颜色
func SetConsoleTextColor(flag int) {
	color := uintptr(flag)
	if mStdColor != color {
		mStdColor = color
		mSetConsoleTextAttributeProc.Call(mStdHandle, mStdColor)
	}
}

// SetFileHide SetFileHide
func SetFileHide(path string, hide bool) error {
	attrs := uint32(syscall.FILE_ATTRIBUTE_NORMAL)
	if hide {
		attrs = syscall.FILE_ATTRIBUTE_HIDDEN
	}
	return syscall.SetFileAttributes(syscall.StringToUTF16Ptr(path), attrs)
}

// SCP SCP
func SCP(remoteAddr, remotePort, remoteUser, remotePass string, srcPath string, dstPath string) (string, error) {
	panic("windows-scp-nonsupport")
}
