package dialog

import (
	"syscall"
	"unsafe"

	"github.com/sirupsen/logrus"
	"leoliu.io/logger"
)

var (
	MsgBoxBtnOk                uint = 0x000000
	MsgBoxBtnOkCancel          uint = 0x000001
	MsgBoxBtnAbortRetryIgnore  uint = 0x000002
	MsgBoxBtnYesNoCancel       uint = 0x000003
	MsgBoxBtnYesNo             uint = 0x000004
	MsgBoxBtnRetryCancel       uint = 0x000005
	MsgBoxBtnCancelTryContinue uint = 0x000006
	MsgBoxIconNone             uint = 0x000000
	MsgBoxIconError            uint = 0x000010
	MsgBoxIconQuestion         uint = 0x000020
	MsgBoxIconWarning          uint = 0x000030
	MsgBoxIconInformation      uint = 0x000040
	MsgBoxDefaultButton1       uint = 0x000000
	MsgBoxDefaultButton2       uint = 0x000100
	MsgBoxDefaultButton3       uint = 0x000200
	MsgBoxDefaultButton4       uint = 0x000300
	MsgBoxTopMost              uint = 0x041000
	MsgBoxService              uint = 0x200000

	MsgBoxSelectOk       = 1
	MsgBoxSelectCancel   = 2
	MsgBoxSelectAbort    = 3
	MsgBoxSelectRetry    = 4
	MsgBoxSelectIgnore   = 5
	MsgBoxSelectYes      = 6
	MsgBoxSelectNo       = 7
	MsgBoxSelectTry      = 10
	MsgBoxSelectContinue = 11
)

// MsgBox create message box
func MsgBox(title string, msg string, flag uint) (int, error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"title":   title,
				"message": msg,
				"flag":    flag,
			}),
		).Debugln("Create message box . . .")
	}

	// Generate message box
	if intLog {
		intLogger.Debugln("Generate message box")
	}
	rtn, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msg))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flag),
	)
	if rtn == 0 {
		rtn, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("GetLastError").Call()
		err := MsgBoxError(uint32(rtn))

		intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
			WithError(err).Errorln("Cannot generate message box")
		return 0, err
	}

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"button_code": rtn,
			}),
		).Debugln("Create message box")
	}
	return int(rtn), nil
}
