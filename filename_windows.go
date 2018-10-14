package dialog

import (
	"unsafe"
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms646839.aspx
type openFileName struct {
	structSize      uint32
	owner           uintptr
	instance        uintptr
	filter          *uint16
	customFilter    *uint16
	maxCustomFilter uint32
	filterIndex     uint32
	file            *uint16
	maxFile         uint32
	fileTitle       *uint16
	maxFileTitle    uint32
	initialDir      *uint16
	title           *uint16
	flags           uint32
	fileOffset      uint16
	fileExtension   uint16
	defExt          *uint16
	custData        uintptr
	fnHook          uintptr
	templateName    *uint16
	pvReserved      unsafe.Pointer
	dwReserved      uint32
	flagsEx         uint32
}

type FileNameFilters map[string][]string

var (
	FileNameShowHidden  uint32 = 0x10000000
	FileNameMultiSelect uint32 = 0x00000200
)
