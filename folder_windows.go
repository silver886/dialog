package dialog

import (
	"errors"
	"reflect"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/sirupsen/logrus"
	"leoliu.io/logger"
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb773205.aspx
type browseInfo struct {
	owner        uintptr
	root         uintptr
	displayName  *uint16
	title        *uint16
	flags        uint32
	callbackFunc uintptr
	lParam       uintptr
	image        int32
}

var (
	FolderEditBox            uint32 = 0x00000010
	FolderNewDialogStyle     uint32 = 0x00000040
	FolderNoNewFolderButton  uint32 = 0x00000200
	FolderBrowseIncludeFiles uint32 = 0x00004000
)

// Folder create Browse For Folder dialog
func Folder(title string, initDir string, flag uint32, exLong bool) (string, error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"title":             title,
				"initial_directory": initDir,
				"flag":              flag,
				"extremely long":    exLong,
			}),
		).Debugln("Create Browse For Folder dialog . . .")
	}
	if initDir == "" {
		// This PC
		initDir = "::{20D04FE0-3AEA-1069-A2D8-08002B30309D}"
	}

	// Set parameters
	bi := &browseInfo{}

	bi.flags = flag
	bi.title, _ = syscall.UTF16PtrFromString(title)

	dir := utf16.Encode([]rune(initDir))
	dirPtr := (*reflect.SliceHeader)(unsafe.Pointer(&dir)).Data
	bi.root, _, _ = syscall.NewLazyDLL("shell32.dll").NewProc("SHSimpleIDListFromPath").Call(dirPtr)

	// Generate Browse For Folder dialog
	if intLog {
		intLogger.Debugln("Generate Browse For Folder dialog . . .")
	}
	rtn, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("SHBrowseForFolderW").Call(
		uintptr(unsafe.Pointer(bi)),
	)
	if rtn == 0 {
		if intLog {
			intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
				Errorln("User cancelled")
		}
		return "", errors.New("User cancelled")
	}

	// Get folder path
	var pathLen int
	if exLong {
		pathLen = 0x7fffffff
	} else {
		pathLen = 0x00000fff
	}
	pathBuf := make([]uint16, pathLen)
	syscall.NewLazyDLL("shell32.dll").NewProc("SHGetPathFromIDListW").Call(
		rtn,
		uintptr(unsafe.Pointer(&pathBuf[0])),
	)
	if intLog {
		intLogger.WithFields(logrus.Fields{}).Debugln("Get folder path")
	}
	return syscall.UTF16ToString(pathBuf), nil
}
