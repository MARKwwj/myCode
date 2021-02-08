package api

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

func init() {
}

// SetConsoleTitle 设置控制台标题
func SetConsoleTitle(title string) {
}

// SetConsoleTextColor 设置控制台文本颜色
func SetConsoleTextColor(flag int) {
	switch flag {
	case 0x000B: //eLogDebugType:
	case 0x000F: //eLogInfoType:
	case 0x000E: //eLogWarnType:
	case 0x000D: //eLogFatalType:
	case 0x000C: //eLogErrorType:
	default:
		// 0x0007 other
	}
}

// SetFileHide SetFileHide
func SetFileHide(path string, hide bool) error {
	return nil
}

// SCP SCP
func SCP(remoteAddr, remotePort, remoteUser, remotePass string, srcPath string, dstPath string) (string, error) {
	srcPath, _ = filepath.Abs(srcPath)
	var out bytes.Buffer
	cmd := exec.Command("sshpass", "-p", remotePass, "scp", "-P", remotePort, "-r", srcPath, fmt.Sprintf("%s@%s:%s", remoteUser, remoteAddr, dstPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Start(); err != nil {
		return out.String(), err
	}
	if err := cmd.Wait(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}
