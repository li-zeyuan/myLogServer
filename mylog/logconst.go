package mylog

const (
	DebugLevel = 0
	TraceLevel = 1
	InfoLevel  = 2
	WarnLevel  = 3
	ErrorLevel = 4
	FatalLevel = 5
)

const (
	LogSplitTypeHour = 0
	LogSplitTypeSize = 1
)

func LevelString(level int) string {
	levelStr := ""
	switch level {
	case DebugLevel:
		levelStr = "DEBUG"
	case TraceLevel:
		levelStr = "TRACE"
	case InfoLevel:
		levelStr = "INFO"
	case WarnLevel:
		levelStr = "WARN"
	case ErrorLevel:
		levelStr = "ERROR"
	case FatalLevel:
		levelStr = "FATAL"
	}
	return levelStr
}