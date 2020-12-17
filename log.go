package main

import (
	"fmt"
	"mylogserver/mylog"
)
/*
对FileLog/ConsoleLog又封装了一层
1、var 一个Logger接口，InitLog根据不用的name创建不同的Log对象
2、Log对象调各自的方法
 */
var log mylog.Logger

func InitLog(name string, config map[string]string) (err error)  {
	switch name {
	case "file":
		log, err = mylog.NewFileLog(config)
	case "console":
		log = mylog.NewConsoleLog()
	default:
		err = fmt.Errorf("unspport log name:%s", name)
	}
	return
}

func Debug(format string, args ...interface{})  {
	log.Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	log.Trace(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatal(format, args...)
}

func Close() {
	log.Close()
}


