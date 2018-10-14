package dialog

import (
	"fmt"
)

type MsgBoxError uint32

func (e MsgBoxError) Error() string {
	return fmt.Sprintf("SystemErrorCode: %#x", e)
}

type FileError int

func (e FileError) Error() string {
	return fmt.Sprintf("CommDlgExtendedError: %#x", e)
}
