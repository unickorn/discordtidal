package discordtidal

import (
	"discordtidal/song"
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

	processIDs = make(map[uint32]interface{})
	timer      = 5

	track  string
	artist string
	st     Status
	cb     = syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		processIDs := GetTidalProcessIDs()

		processId := GetWindowThreadProcessId(uintptr(h))
		// skip if process id is not one of the TIDAL.exe ones
		_, ok := processIDs[uint32(processId)]
		if !ok {
			return 1
		}

		title := GetWindowText(h)
		if !strings.Contains(title, " - ") || strings.Contains(title, "{") {
			if song.Current != nil {
				st = Paused
			} else {
				st = Opened
			}
			return 1
		}

		s := strings.Split(title, " - ")
		track = strings.Join(s[:len(s)-1], " - ")
		artist = s[len(s)-1] // just assume it's always the song that has dashes lmao
		st = Playing
		return 0
	})
)

// GetSong tries to get the song and artist from Tidal window title.
func GetSong() (string, string, Status) {
	EnumWindows(cb, 0)
	return track, artist, st
}

// GetTidalProcessIDs returns a map[TIDAL.exe pid]0 because we will abuse this
// map to check if a process ID is one of TIDAL ones later on.
func GetTidalProcessIDs() map[uint32]interface{} {
	timer--
	if len(processIDs) > 0 && timer > 0 {
		return processIDs
	}
	timer = 5
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
	for {
		if parseProcessName(procEntry.ExeFile) == "TIDAL.exe" {
			processIDs[procEntry.ProcessID] = struct{}{}
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
// Copies the text of the specified window's title bar (if it has one) into a buffer.
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
