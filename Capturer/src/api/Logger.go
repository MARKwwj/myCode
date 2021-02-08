package api

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	mFilterLevel int
	mWaitGroup   sync.WaitGroup
	mLogFilePath string
)

// StartService 启动日志服务
func StartService() error {
	mWaitGroup.Add(2)
	EnableDebugInfo(true)
	go goroutinePrintLog()
	go goroutineWriteLog()
	return nil
}

// CloseService 关闭日志服务
func CloseService() {
	ProtectError()
	mLogQueue <- nil
	mWaitGroup.Wait()
}

// EnableDebugInfo 启用Debug日志信息
func EnableDebugInfo(enable bool) {
	if enable {
		mFilterLevel = 0
	} else {
		mFilterLevel = eLogInfoType
	}
}

// ProtectError ProtectError
func ProtectError() {
	if err := recover(); err != nil {
		buff := make([]byte, 1024*4)
		n := runtime.Stack(buff, false)
		Error("%v · Call stack:\n%s", err, buff[:n])
	}
}

// Debug Debug
func Debug(format string, a ...interface{}) {
	if IsDebug() && mFilterLevel <= eLogDebugType {
		postLoginfo(eLogDebugType, format, a...)
	}
}

// Info Info
func Info(format string, a ...interface{}) {
	if mFilterLevel <= eLogInfoType {
		postLoginfo(eLogInfoType, format, a...)
	}
}

// Warn Warn
func Warn(format string, a ...interface{}) {
	if mFilterLevel <= eLogWarnType {
		postLoginfo(eLogWarnType, format, a...)
	}
}

// Fatal Fatal
func Fatal(format string, a ...interface{}) {
	if mFilterLevel <= eLogFatalType {
		postLoginfo(eLogFatalType, format, a...)
	}
}

// Error Error
func Error(format string, a ...interface{}) {
	if mFilterLevel <= eLogErrorType {
		postLoginfo(eLogErrorType, format, a...)
	}
}

const (
	eLogDebugType = 0x01
	eLogInfoType  = 0x02
	eLogWarnType  = 0x04
	eLogFatalType = 0x08
	eLogErrorType = 0x10
)

func getLogTypeName(eType int) string {
	switch eType {
	case eLogDebugType:
		return "[D]"
	case eLogInfoType:
		return "[I]"
	case eLogWarnType:
		return "[W]"
	case eLogFatalType:
		return "[F]"
	case eLogErrorType:
		return "[E]"
	}
	return "[U]"
}

func postLoginfo(logType int, format string, a ...interface{}) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	nowTime := time.Now()
	fmt.Fprintf(buffer, "[%04d-%02d-%02d %02d:%02d:%02d]%s",
		nowTime.Year(), nowTime.Month(), nowTime.Day(),
		nowTime.Hour(), nowTime.Minute(), nowTime.Second(), getLogTypeName(logType))
	if a == nil {
		fmt.Fprintf(buffer, "%s", format)
	} else {
		fmt.Fprintf(buffer, format, a...)
	}
	mLogQueue <- &stLogInfo{
		mType: logType,
		mInfo: string(buffer.Bytes()),
	}
}

type stLogInfo struct {
	mType int
	mInfo string
}

var (
	mLogQueue     = make(chan *stLogInfo, 64)
	mLogFileQueue = make(chan *stLogInfo, 64)
)

func goroutinePrintLog() {
	defer func() {
		mLogFileQueue <- nil
		mWaitGroup.Done()
	}()
	for info := range mLogQueue {
		if info == nil {
			return
		}
		switch info.mType {
		case eLogDebugType:
			SetConsoleTextColor(0x00F9)
		case eLogInfoType:
			SetConsoleTextColor(0x000F)
		case eLogWarnType:
			SetConsoleTextColor(0x000E)
		case eLogFatalType:
			SetConsoleTextColor(0x000D)
		case eLogErrorType:
			SetConsoleTextColor(0x000C)
		default:
			SetConsoleTextColor(0x0007)
		}
		fmt.Fprintln(os.Stdout, info.mInfo)
		mLogFileQueue <- info
	}
}

func goroutineWriteLog() {
	defer mWaitGroup.Done()
	for info := range mLogFileQueue {
		if info == nil {
			return
		}
		for {
			file, err := os.OpenFile(getLogFileName(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				Fatal("Write to log file failed: %s", err.Error())
				time.Sleep(time.Second)
				continue
			}
			file.Write([]byte(info.mInfo))
			file.Write([]byte("\r\n"))
			file.Close()
			break
		}
	}
}

func getLogFileName() string {
	if len(mLogFilePath) <= 0 {
		fileName := filepath.Base(os.Args[0])
		fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
		mLogFilePath = fileName + ".log"
	}
	return mLogFilePath
}
