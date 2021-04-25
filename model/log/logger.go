package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// 默认日志等级为debug
var currentLevel = LevelDebug

// 日志等级
type logLevel int

// 日志等级
const (
	LevelErr logLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

func Debug(msg string) {
	if enable(LevelDebug) {
		fmt.Printf(getInfo(LevelDebug, msg))
	}
}

func Info(msg string) {
	if enable(LevelInfo) {
		fmt.Printf(getInfo(LevelInfo, msg))
	}
}

func Warn(msg string) {
	if enable(LevelWarn) {
		fmt.Printf(getInfo(LevelWarn, msg))
	}
}

func Error(msg string) {
	if enable(LevelErr) {
		fmt.Printf(getInfo(LevelErr, msg))
		os.Exit(1)
	}
}

// Fatal 输出错误并退出
func Fatal(msg string) {
	fmt.Printf(getInfo(LevelErr, msg))
	os.Exit(1)
}

// SetLevel 设置日志等级
func SetLevel(level logLevel) {
	switch level {
	case LevelErr, LevelWarn, LevelInfo, LevelDebug:
		currentLevel = level
	default:
		Error("没有这个日志等级")
	}
}

// 是否启动该日志等级
func enable(level logLevel) bool {
	return level <= currentLevel
}

// 获取日志等级字符串
func getLevel(level logLevel) (l string) {
	switch level {
	case LevelDebug:
		l = "DEBUG"
	case LevelInfo:
		l = Green("INFO ")
	case LevelWarn:
		l = Yellow("WARN ")
	case LevelErr:
		l = Red("ERROR")
	}
	return
}

// 获取调用方法信息
func getInfo(level logLevel, msg string) string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		log.Fatal("日志获取信息错误")
	}
	funName := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	now := time.Now().Format("2006-01-02 15:04:05")
	funInfo := Cyan(fmt.Sprintf("%s.%s:%s:%d", funName[0], path.Base(file), funName[1], line))
	return fmt.Sprintf("%s %s %s : %s\n", now, getLevel(level), funInfo, msg)
}
