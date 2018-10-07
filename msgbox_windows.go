package dialog

import (
	"github.com/sirupsen/logrus"
	"leoliu.io/execute"
	"leoliu.io/logger"
)

// MsgBox create message box
func MsgBox(title string, msg string, args ...string) (output string, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"title":     title,
				"message":   msg,
				"arguments": args,
			}),
		).Debugln("Create message box . . .")
	}

	// Create message box
	cmd, err := BgMsgBox(title, msg, args...)
	if err != nil {
		if intLog {
			intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
				WithError(err).Errorln("Cannot create message box")
		}
		return
	}
	cmd.Wait()
	output = cmd.Strout()

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"output": output,
			}),
		).Debugln("Create message box")
	}
	return
}

// BgMsgBox create message box in the background
func BgMsgBox(title string, msg string, args ...string) (cmd *execute.Cmd, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"title":     title,
				"message":   msg,
				"arguments": args,
			}),
		).Debugln("Create message box in the background . . .")
	}

	// Parse arguments
	btn := "0x0"
	icon := "0x00"
	defaultBtn := "0x00"
	topMost := "false"
	for _, val := range args {
		switch val {
		// Set button
		case "Ok":
			btn = "0x0"
		case "O":
			btn = "0x0"
		case "OkCancel":
			btn = "0x1"
		case "OC":
			btn = "0x1"
		case "AbortRetryIgnore":
			btn = "0x2"
		case "ARI":
			btn = "0x2"
		case "YesNoCancel":
			btn = "0x3"
		case "YNC":
			btn = "0x3"
		case "YesNo":
			btn = "0x4"
		case "YN":
			btn = "0x4"
		case "RetryCancel":
			btn = "0x5"
		case "RC":
			btn = "0x5"

		// Set icon
		case "None":
			icon = "0x00"
		case "N":
			icon = "0x00"
		case "Error":
			icon = "0x10"
		case "E":
			icon = "0x10"
		case "Question":
			icon = "0x20"
		case "Q":
			icon = "0x20"
		case "Warning":
			icon = "0x30"
		case "W":
			icon = "0x30"
		case "Information":
			icon = "0x40"
		case "I":
			icon = "0x40"

		// Set default button
		case "DefaultButton1":
			defaultBtn = "0x000"
		case "DB1":
			defaultBtn = "0x000"
		case "DefaultButton2":
			defaultBtn = "0x100"
		case "DB2":
			defaultBtn = "0x100"
		case "DefaultButton3":
			defaultBtn = "0x200"
		case "DB3":
			defaultBtn = "0x200"

		// Set top most
		case "TopMost":
			topMost = "true"
		case "TM":
			topMost = "true"
		}
	}
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"button":         btn,
				"icon":           icon,
				"default_button": defaultBtn,
				"top_most":       topMost,
			}),
		).Debugln("Message box arguments")
	}

	// Generate command
	command := []string{"[void] [System.Reflection.Assembly]::LoadWithPartialName('System.Windows.Forms')",
		"; $FrmMain = New-Object 'System.Windows.Forms.Form'",
		"; $FrmMain.TopMost = $" + topMost,
		"; $Answer = [System.Windows.Forms.MessageBox]::Show($FrmMain, '" + msg + "', '" + title + "', " + btn + ", " + icon + ", " + defaultBtn + ")",
		"; $FrmMain.Close()",
		"; $FrmMain.Dispose()",
		"; Write-Output $Answer",
	}
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"command": command,
			}),
		).Debugln("Generate command")
	}

	// Create message box in the background
	cmd, err = execute.Start(
		true,
		"powershell", command...,
	)
	if err != nil {
		intLogger.WithFields(logger.DebugInfo(1, logrus.Fields{})).
			WithError(err).Errorln("Cannot create message box")
		return nil, err
	}

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"command_object": cmd,
			}),
		).Debugln("Create message box in the background")
	}
	return cmd, nil
}
