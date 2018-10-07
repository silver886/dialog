package dialog

import (
	"leoliu.io/logger"
)

var (
	intLog    bool
	intLogger *logger.Entry
)

// SetLogger set internal logger for logging
func SetLogger(extLogger *logger.Logger) {
	intLogger = extLogger.WithField("prefix", "file")
	intLog = true
}

// ResetLogger reset internal logger
func ResetLogger() {
	intLogger = nil
	intLog = false
}
