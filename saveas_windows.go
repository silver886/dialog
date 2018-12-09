package dialog

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// GetNewFileName create new file name dialog
func GetNewFileName(title string, initDir string, filter FileNameFilters, flag uint32, exLong bool) (string, error) {
	// Set parameters
	ofn := &openFileName{}
	ofn.structSize = uint32(unsafe.Sizeof(*ofn))

	ofn.flags = 0x02080002 | flag
	ofn.title, _ = syscall.UTF16PtrFromString(title)
	ofn.initialDir, _ = syscall.UTF16PtrFromString(initDir)

	var fileLen int
	if exLong {
		fileLen = 0x7fffffff
	} else {
		fileLen = 0x00000fff
	}
	fileBuf := make([]uint16, fileLen)

	ofn.file = utf16ptr(fileBuf)
	ofn.maxFile = uint32(fileLen)

	var filters []uint16
	var filtersStr []string
	for desc, exts := range filter {
		// "Music File\0*.mp3;*.ogg;*.wav;\0"
		filters = append(filters, utf16.Encode([]rune(desc))...)
		filters = append(filters, 0)
		for _, ext := range exts {
			s := fmt.Sprintf("*.%s;", ext)
			filters = append(filters, utf16.Encode([]rune(s))...)
		}
		filters = append(filters, 0)
		filtersStr = append(filtersStr, exts[0])
	}
	if filters != nil {
		// Two extra NUL chars to terminate the list
		filters = append(filters, 0, 0)
		ofn.filter = utf16ptr(filters)
	}

	// Generate new file name dialog
	rtn, _, _ := syscall.NewLazyDLL("comdlg32.dll").NewProc("GetSaveFileNameW").Call(
		uintptr(unsafe.Pointer(ofn)),
	)
	if rtn == 0 {
		rtn, _, _ := syscall.NewLazyDLL("comdlg32.dll").NewProc("CommDlgExtendedError").Call()
		if uint32(rtn) == 0 {
			return "", errors.New("User cancelled")
		}
		err := FileError(uint32(rtn))

		return "", err
	}

	// Get new file name
	var fileName strings.Builder
	i := 0
	for i < fileLen && (fileBuf[i] != 0 || fileBuf[i+1] != 0) {
		i++
	}
	fileName.WriteString(string(utf16.Decode(fileBuf[:i])))

	// Get file extension
	if ofn.filterIndex != 0 {
		fileName.WriteString(".")
		fileName.WriteString(filtersStr[ofn.filterIndex-1])
	}

	return fileName.String(), nil
}
