package mylog

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLog struct {
	logPath       string
	logName       string
	file          *os.File
	warnFile      *os.File
	logDataChan   chan *LogData
	logSplitType  int
	logSplitSize  int64
	lastSplitHour int
}

// 构造FileLog对象，返回interface
func NewFileLog(config map[string]string) (Logger, error) {
	logPath, ok := config["log_path"]
	if !ok {
		err := fmt.Errorf("not found log_path")
		return nil, err
	}
	logName, ok := config["log_name"]
	if !ok {
		err := fmt.Errorf("not found log_name")
		return nil, err
	}
	logChanSize, ok := config["log_chan_size"]
	if !ok {
		logChanSize = "50000"
	}
	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 50000
	}

	var logSplitType = LogSplitTypeSize
	var logSplitSize int64
	splitType, ok := config["log_split_type"]
	if !ok {
		splitType = "hour"
	} else {
		if splitType == "size" {
			splitSize, ok := config["log_split_size"]
			if !ok {
				splitSize = "104857600"
			}
			logSplitSize, err = strconv.ParseInt(splitSize, 10, 64)
			if err != nil {
				logSplitSize = 104857600
			}
			logSplitType = LogSplitTypeSize
		} else {
			logSplitType = LogSplitTypeHour
		}
	}

	logFile := new(FileLog)
	logFile.logPath = logPath
	logFile.logName = logName
	logFile.logDataChan = make(chan *LogData, chanSize)
	logFile.logSplitType = logSplitType
	logFile.logSplitSize = logSplitSize
	logFile.lastSplitHour = time.Now().Hour()

	logFile.Init()

	return logFile, nil
}

// 创建日志文件
func (f *FileLog) Init() {
	// 普通日志
	fileName := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed, err: %v", fileName, err))
	}
	f.file = file

	// 错误日志
	warnFileName := fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	warnFile, err := os.OpenFile(warnFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755) // os.O_CREATE 创建文件 os.O_APPEND 追加写入 os.O_WRONLY 只写操作
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed, err: %v", warnFileName, err))
	}
	f.warnFile = warnFile
}

func (f *FileLog) writeLogBackGround() {
	for logData := range f.logDataChan {
		var file = f.file
		if logData.IsWarn {
			file = f.warnFile
		}

		fmt.Fprintf(file, "%s %s [%s/%s:%d] %s\n", logData.TimeStr, logData.LevelStr,
			logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
	}
}

func (f *FileLog) checkSplitFile(isWarn bool) {
	if f.logSplitType == LogSplitTypeHour {
		f.splitHour(isWarn)
		return
	}
	f.splitSize(isWarn)
}

// 根据文件大小分割
func (f *FileLog) splitSize(isWarn bool) {
	file := f.file
	defer file.Close()
	if isWarn {
		file = f.warnFile
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return
	}
	fileSize := fileInfo.Size()
	if fileSize <= f.logSplitSize {
		return
	}

	backupFileName := ""
	fileName := ""
	now := time.Now()
	if isWarn {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%s", f.logPath, f.logName, now.Format("20060102150405"))
		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%s", f.logPath, f.logName, now.Format("20060102150405"))
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}
	os.Rename(fileName, backupFileName)

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}

	if isWarn {
		f.warnFile = file
	} else {
		f.file = file
	}
}

// 根据时间来进行分割
func (f *FileLog) splitHour(isWarn bool) {

	now := time.Now()
	hour := now.Hour()

	if hour == f.lastSplitHour {
		return
	}

	f.lastSplitHour = hour

	var backupFileName string
	var fileName string

	if isWarn {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%s", f.logPath, f.logName, now.Format("20060102150405"))
		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%s", f.logPath, f.logName, now.Format("20060102150405"))
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}

	file := f.file
	if isWarn {
		file = f.warnFile
	}
	file.Close()
	os.Rename(fileName, backupFileName)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}

	if isWarn {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLog) Debug(format string, args ...interface{}) {
	logData := writeLog(DebugLevel, format, args...)
	select {
	case f.logDataChan <- logData:
	default:

	}
}

func (f *FileLog) Trace(format string, args ...interface{}) {
	writeLog(TraceLevel, format, args...)
}

func (f *FileLog) Info(format string, args ...interface{}) {
	logData := writeLog(InfoLevel, format, args)
	select {
	case f.logDataChan <- logData:
		go f.writeLogBackGround()
	default:
	}
}

func (f *FileLog) Warn(format string, args ...interface{}) {
	logData := writeLog(WarnLevel, format, args)
	select {
	case f.logDataChan <- logData:
		go f.writeLogBackGround()
	default:
	}
	time.Sleep(time.Second * 1)
}

func (f *FileLog) Error(format string, args ...interface{}) {
	fmt.Fprintf(f.warnFile, format, args...)
	fmt.Fprintln(f.warnFile)
}

func (f *FileLog) Fatal(format string, args ...interface{}) {
	fmt.Fprintf(f.warnFile, format, args...)
	fmt.Fprintln(f.warnFile)
}

func (f *FileLog) Close() {
	f.file.Close()
	f.warnFile.Close()
}
