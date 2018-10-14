package dialog

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/sirupsen/logrus"
	"leoliu.io/logger"
)

// SaveAs create Save As dialog
func SaveAs(title string, initDir string, filter FileNameFilters, flag uint32, exLong bool) (string, error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"title":             title,
				"initial_directory": initDir,
				"filter":            filter,
				"flag":              flag,
				"extremely long":    exLong,
			}),
		).Debugln("Create Save As dialog . . .")
	}

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

	// Generate Save As dialog
	if intLog {
		intLogger.Debugln("Generate Save As dialog . . .")
	}
	rtn, _, _ := syscall.NewLazyDLL("comdlg32.dll").NewProc("GetSaveFileNameW").Call(uintptr(unsafe.Pointer(ofn)))
	if rtn == 0 {
		rtn, _, _ := syscall.NewLazyDLL("comdlg32.dll").NewProc("CommDlgExtendedError").Call()
		if uint32(rtn) == 0 {
			return "", errors.New("User cancelled")
		}
		err := FileError(uint32(rtn))

		if intLog {
			intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
				WithError(err).Errorln("Cannot generate Save As dialog")
		}
		return "", err
	}

	// Get file name
	var fileName strings.Builder
	i := 0
	for i < fileLen && (fileBuf[i] != 0 || fileBuf[i+1] != 0) {
		i++
	}
	fileName.WriteString(string(utf16.Decode(fileBuf[:i])))
	if intLog {
		intLogger.WithFields(logrus.Fields{
			"file_name": fileName.String(),
		}).Debugln("Get file name")
	}

	// Get file extension
	if ofn.filterIndex != 0 {
		fileName.WriteString(".")
		fileName.WriteString(filtersStr[ofn.filterIndex-1])
		if intLog {
			intLogger.WithFields(logrus.Fields{
				"file_extension": filtersStr[ofn.filterIndex-1],
			}).Debugln("Get file extension")
		}
	}

	return fileName.String(), nil
}
