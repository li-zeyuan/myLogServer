package mylog

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

type LogData struct {
	Message  string
	TimeStr  string
	LevelStr string
	FileName string
	FuncName string
	LineNo   int
	IsWarn   bool
}



func GetLineInfo() (string, string, int) {
	fileName, funcName, lineNo := "", "", 0
	pc, file, line, ok := runtime.Caller(0)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}

	return fileName, funcName, lineNo
}

func writeLog(level int, format string, args ...interface{}) *LogData {
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")
	levelStr := LevelString(level)
	fileName, funcName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)
	isWarn := level >= WarnLevel && level <= FatalLevel

	logData := new(LogData)
	logData.Message = msg
	logData.TimeStr = nowStr
	logData.LevelStr = levelStr
	logData.FileName = fileName
	logData.FuncName = funcName
	logData.LineNo = lineNo
	logData.IsWarn = isWarn

	return logData
}
