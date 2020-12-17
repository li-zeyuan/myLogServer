package mylog

import (
	"fmt"
	"os"
)

type ConsoleLog struct {

}

func NewConsoleLog() Logger {
	logCon := &ConsoleLog{}
	return logCon
}

func (c *ConsoleLog) Debug(format string, args ...interface{}) {
	logData := writeLog(DebugLevel, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s [%s/%s:%d] %s\n", logData.TimeStr,
		logData.LevelStr, logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLog) Trace(format string, args ...interface{}) {

}

func (c *ConsoleLog) Info(format string, args ...interface{}) {

}

func (c *ConsoleLog) Warn(format string, args ...interface{}) {
	logData := writeLog(WarnLevel, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s [%s/%s:%d] %s\n", logData.TimeStr,
		logData.LevelStr, logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLog) Error(format string, args ...interface{}) {
}

func (c *ConsoleLog) Fatal(format string, args ...interface{}) {

}

func (c *ConsoleLog) Close() {

}