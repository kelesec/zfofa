package filelog

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var logger *log.Logger
var fileLoggerOnce sync.Once // 控制日志文件只打开一次，减少性能消耗
var FileContext context.Context

func init() {
	// 日志文件名
	logFile := time.Now().Format("2006/01/02") + ".log"
	logFile = strings.ReplaceAll(logFile, "/", "-")
	logFile = strings.ReplaceAll(logFile, " ", "-")
	logFile = strings.ReplaceAll(logFile, ":", "-")

	fileLoggerOnce.Do(func() {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger = log.New(file, "", 0)
		}
	})
}

// save 保存日志到文件中
func save(msg, level string) {
	currentTime := time.Now().Format("2006/01/02 15:04:05")
	fileMsg := fmt.Sprintf("[%s] [%s] %s", currentTime, level, msg)

	if logger != nil {
		logger.Printf(fileMsg)
	}

	if FileContext != nil {
		fileMsg = strings.TrimSpace(fileMsg)
		runtime.EventsEmit(FileContext, "PushLog", fileMsg) // 触发前端响应日志事件
	}
}

func Info(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	save(msg, "INFO")
}

func Debug(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	save(msg, "DEBUG")
}

func Warning(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	save(msg, "WARNING")
}

func Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	save(msg, "ERROR")
}

func Fatalf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	save(msg, "FATAL")
	os.Exit(1)
}
