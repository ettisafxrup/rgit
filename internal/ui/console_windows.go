//go:build windows

package ui

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	enableVirtualTerminalProcessing = 0x0004
	utf8CodePage                    = 65001
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode     = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode     = kernel32.NewProc("SetConsoleMode")
	procSetConsoleOutputCP = kernel32.NewProc("SetConsoleOutputCP")
	procSetConsoleCP       = kernel32.NewProc("SetConsoleCP")

	procGetFileInfoByHandle = kernel32.NewProc("GetFileInformationByHandleEx")
)

// PrepareConsole switches the Windows console to UTF-8 so that symbols such
// as ✓ and tree lines render correctly in cmd.exe and PowerShell.
func PrepareConsole() {
	procSetConsoleOutputCP.Call(utf8CodePage)
	procSetConsoleCP.Call(utf8CodePage)
}

func consoleMode(f *os.File) (uint32, bool) {
	var mode uint32
	ok, _, _ := procGetConsoleMode.Call(f.Fd(), uintptr(unsafe.Pointer(&mode)))
	return mode, ok != 0
}

func isTerminal(f *os.File) bool {
	if _, ok := consoleMode(f); ok {
		return true
	}
	return isMinttyPipe(f)
}

// isMinttyPipe detects Git Bash / Cygwin terminals (mintty). They expose a
// named pipe instead of a console, called like
// \msys-<id>-pty<n>-to-master, and understand ANSI colors.
func isMinttyPipe(f *os.File) bool {
	const fileNameInfo = 2 // FILE_INFO_BY_HANDLE_CLASS: FileNameInfo
	buf := make([]uint16, 1024)
	ret, _, _ := procGetFileInfoByHandle.Call(f.Fd(), fileNameInfo,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)*2))
	if ret == 0 {
		return false
	}
	// The buffer starts with the name length in bytes, then the UTF-16 name.
	length := int(uint32(buf[0])|uint32(buf[1])<<16) / 2
	if length <= 0 || length > len(buf)-2 {
		return false
	}
	name := syscall.UTF16ToString(buf[2 : 2+length])
	isMsys := strings.Contains(name, "msys-") || strings.Contains(name, "cygwin-")
	return isMsys && strings.Contains(name, "-pty") && strings.HasSuffix(name, "-master")
}

// enableVirtualTerminal asks the console to interpret ANSI escape sequences.
func enableVirtualTerminal(f *os.File) bool {
	mode, ok := consoleMode(f)
	if !ok {
		return true // not a console (e.g. mintty), which handles ANSI itself
	}
	if mode&enableVirtualTerminalProcessing != 0 {
		return true
	}
	ret, _, _ := procSetConsoleMode.Call(f.Fd(), uintptr(mode|enableVirtualTerminalProcessing))
	return ret != 0
}
