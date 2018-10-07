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
				"filter":            filter,
				"initial_directory": initDir,
			}),
		).Debugln("Create Save As dialog . . .")
	}

	// Create Save As dialog
	cmd, err := BgSaveAs(filter, initDir)
	if err != nil {
		if intLog {
			intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
				WithError(err).Errorln("Cannot create Save As dialog")
		}
		return
	}
	cmd.Wait()

	// Parse output
	result := strings.Split(cmd.Strout(), "\r\n")
	if result[0] != "OK" {
		return "", errors.New("Cancelled by user")
	}
	path = result[1]

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"path": path,
			}),
		).Debugln("Create Save As dialog")
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
		).Debugln("Create Save As dialog in the background . . .")
	}

	// Parse arguments
	initDir = strings.Replace(initDir, "/", "\\", -1)

	// Generate command
	command := []string{"[void] [System.Reflection.Assembly]::LoadWithPartialName('System.Windows.Forms')",
		"; $SaveFile = New-Object System.Windows.Forms.SaveFileDialog -Property @{",
		"Filter = '" + filter + "'",
		"; InitialDirectory = '" + initDir + "'",
		"}",
		"; Write-Output $SaveFile.ShowDialog()",
		"; Write-Output $SaveFile.FileName",
		"; $SaveFile.Dispose()",
	}
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"command": command,
			}),
		).Debugln("Generate command")
	}

	// Create Save As dialog in the background
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
		).Debugln("Create Save As dialog in the background")
	}

	if err != nil {
		return nil, err
	}
	return cmd, nil
}
