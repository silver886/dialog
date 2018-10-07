package dialog

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"leoliu.io/execute"
	"leoliu.io/logger"
)

// SaveAs create Save As dialog
func SaveAs(filter string, initDir string) (path string, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"filter": filter,
				"path":   path,
			}),
		).Debugln("Create save as dialog . . .")
	}

	cmd, err := BgSaveAs(filter, initDir)
	if err != nil {
		if intLog {
			intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
				WithError(err).Errorln("Cannot create message box . . .")
		}
		return
	}
	cmd.Wait()
	path = cmd.Strout()

	if path == "" {
		return "", errors.New("Cancelled by user")
	}

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"path": path,
			}),
		).Debugln("Create save as dialog")
	}
	return
}

// BgSaveAs create Save As dialog in the background
func BgSaveAs(filter string, initDir string) (cmd *execute.Cmd, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"filter":            filter,
				"initial_directory": initDir,
			}),
		).Debugln("Create save as dialog . . .")
	}

	initDir = strings.Replace(initDir, "/", "\\", -1)

	// Generate command
	command := []string{"[void] [System.Reflection.Assembly]::LoadWithPartialName('System.Windows.Forms')",
		"; $SaveChooser = New-Object System.Windows.Forms.SaveFileDialog -Property @{",
		"Filter = '" + filter + "'",
		"; InitialDirectory = '" + initDir + "'",
		"}",
		"; $SaveChooser.ShowDialog() | Out-Null",
		"; Write-Output $SaveChooser.FileName",
	}
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"command": command,
			}),
		).Debugln("Generate command")
	}

	// Generate message box
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
		).Debugln("Generate message box")
	}

	if err != nil {
		return nil, err
	}
	return cmd, nil
}
