package main

import (
	"fmt"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	user32                       = syscall.MustLoadDLL("user32.dll")
	procGetWindowTextW           = user32.MustFindProc("GetWindowTextW")
	procEnumWindows              = user32.MustFindProc("EnumWindows")
	procGetWindowThreadProcessId = user32.MustFindProc("GetWindowThreadProcessId")
)

func main() {
	song, artist := GetSong()
	fmt.Printf("You are listening to: %v by %v", song, artist)
}

// GetSong tries to get the song and artist from Tidal's window title.
func GetSong() (string, string) {
	processIDs := GetTidalProcessIDs()

	// Get window name that has the same process id
	var song string
	var artist string
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		processId := GetWindowThreadProcessId(uintptr(h))
		// skip if process id is not one of the TIDAL.exe ones
		_, ok := processIDs[uint32(processId)]
		if !ok {
			return 1
		}

		title := GetWindowText(h)
		if !strings.Contains(title, " - ") || strings.Contains(title, "}") {
			return 1
		}

		s := strings.Split(title, " - ")
		song = s[0]
		artist = s[1]
		return 0
	})
	EnumWindows(cb, 0)

	return song, artist
}

// GetTidalProcessIDs returns a map[TIDAL.exe pid]0 because we will abuse this
// map to check if a process ID is one of TIDAL's later on.
func GetTidalProcessIDs() map[uint32]byte {
	// Get a snapshot of all processes
	snapshot, err := syscall.CreateToolhelp32Snapshot(0x00000002, 0) // the flag is TH32CS_SNAPPROCESS
	if err != nil {
		panic(err)
	}
	defer func(handle syscall.Handle) {
		if err := syscall.CloseHandle(handle); err != nil {
			panic(err)
		}
	}(snapshot)

	var procEntry syscall.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err = syscall.Process32First(snapshot, &procEntry); err != nil {
		panic(err)
	}

	// Get processes with TIDAL.exe exefile
	processIDs := make(map[uint32]byte)
	for {
		if parseProcessName(procEntry.ExeFile) == "TIDAL.exe" {
			processIDs[procEntry.ProcessID] = 0
		}
		if err = syscall.Process32Next(snapshot, &procEntry); err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			panic(err)
		}
	}

	return processIDs
}

// GetWindowText wraps around user32.GetWindowTextW, which apparently gets
// window text by a window handle.
func GetWindowText(hwnd syscall.Handle) string {
	b := make([]uint16, 200)
	_, _, _ = procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&b[0])), uintptr(int32(len(b))))
	return syscall.UTF16ToString(b)
}

// GetWindowThreadProcessId wraps around user32.GetWindowThreadProcessId. Below is the API doc from Microsoft.
// Copies the text of the specified window's title bar (if it has one) into a buffer. If the specified window is a
// control, the text of the control is copied. However, GetWindowText cannot retrieve the text of a control
// in another application.
func GetWindowThreadProcessId(hwnd uintptr) uintptr {
	var processId uintptr = 0
	_, _, _ = procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processId)))
	return processId
}

// EnumWindows wraps around user32.EnumWindows. Below is the API doc from Microsoft.
// Enumerates all top-level windows on the screen by passing the handle to each window, in turn, to an
// application-defined callback function. EnumWindows continues until the last top-level window is enumerated or
// the callback function returns FALSE.
func EnumWindows(enumFunc uintptr, lparam uintptr) {
	_, _, _ = procEnumWindows.Call(enumFunc, lparam)
}

// parseProcessName parses the whatever-that-is exeFile name into a string.
func parseProcessName(exeFile [syscall.MAX_PATH]uint16) string {
	for i, v := range exeFile {
		if v <= 0 {
			return string(utf16.Decode(exeFile[:i]))
		}
	}
	return ""
}
