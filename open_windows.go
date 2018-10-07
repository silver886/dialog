package dialog

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"leoliu.io/execute"
	"leoliu.io/logger"
)

// Open create Open dialog
func Open(multi bool, filter string, initDir string) (path []string, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"multi":             multi,
				"filter":            filter,
				"initial_directory": initDir,
			}),
		).Debugln("Create Open dialog . . .")
	}

	// Create Open dialog
	cmd, err := BgOpen(multi, filter, initDir)
	if err != nil {
		if intLog {
			intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
				WithError(err).Errorln("Cannot create Open dialog")
		}
		return
	}
	cmd.Wait()
	rawPath := cmd.Strout()

	// Parse output
	path = strings.Split(rawPath, "\r\n")
	if path[0] == "" {
		return nil, errors.New("Cancelled by user")
	}

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"path": path,
			}),
		).Debugln("Create Open dialog")
	}
	return
}

// BgOpen create Open dialog in the background
func BgOpen(multi bool, filter string, initDir string) (cmd *execute.Cmd, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"multi":             multi,
				"filter":            filter,
				"initial_directory": initDir,
			}),
		).Debugln("Create Open dialog in the background . . .")
	}

	// Parse arguments
	multiStr := "false"
	if multi {
		multiStr = "true"
	}
	initDir = strings.Replace(initDir, "/", "\\", -1)

	// Generate command
	command := []string{"[void] [System.Reflection.Assembly]::LoadWithPartialName('System.Windows.Forms')",
		"; $FileBrowser = New-Object System.Windows.Forms.OpenFileDialog -Property @{",
		"Multiselect = $" + multiStr,
		"; Filter = '" + filter + "'",
		"; InitialDirectory = '" + initDir + "'",
		"}",
		"; $FileBrowser.ShowDialog() | Out-Null",
		"; Write-Output $FileBrowser.FileNames",
	}
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"command": command,
			}),
		).Debugln("Generate command")
	}

	// Create Open dialog in the background
	cmd, err = execute.Start(
		true,
		"powershell", command...,
	)
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"command_object": cmd,
				"error":          err,
			}),
		).Debugln("Create Open dialog in the background")
	}

	if err != nil {
		return nil, err
	}
	return cmd, nil
}
