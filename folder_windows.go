package dialog

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"leoliu.io/execute"
	"leoliu.io/logger"
)

// Folder create Browse For Folder dialog
func Folder(msg string, newFolder bool, initDir string) (path string, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"message":           msg,
				"new_folder":        newFolder,
				"initial_directory": initDir,
			}),
		).Debugln("Create Browse For Folder dialog . . .")
	}

	// Create Browse For Folder dialog
	cmd, err := BgFolder(msg, newFolder, initDir)
	if err != nil {
		if intLog {
			intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
				WithError(err).Errorln("Cannot create Browse For Folder dialog")
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
		).Debugln("Create Browse For Folder dialog")
	}
	return
}

// BgFolder create Browse For Folder dialog in the background
func BgFolder(msg string, newFolder bool, initDir string) (cmd *execute.Cmd, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"message":           msg,
				"new_folder":        newFolder,
				"initial_directory": initDir,
			}),
		).Debugln("Create Browse For Folder dialog in the background . . .")
	}

	// Parse arguments
	newFolderStr := "false"
	if newFolder {
		newFolderStr = "true"
	}
	initDir = strings.Replace(initDir, "/", "\\", -1)

	// Generate command
	command := []string{"[void] [System.Reflection.Assembly]::LoadWithPartialName('System.Windows.Forms')",
		"; $FileBrowser = New-Object System.Windows.Forms.FolderBrowserDialog -Property @{",
		"Description = '" + msg + "'",
		"; ShowNewFolderButton = $" + newFolderStr,
		"; SelectedPath = '" + initDir + "'",
		"; RootFolder='MyComputer'",
		"}",
		"; $FileBrowser.ShowDialog()",
		"; Write-Output $FileBrowser.SelectedPath",
	}
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"command": command,
			}),
		).Debugln("Generate command")
	}

	// Create Browse For Folder dialog in the background
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
		).Debugln("Create Browse For Folder dialog in the background")
	}

	if err != nil {
		return nil, err
	}
	return cmd, nil
}
