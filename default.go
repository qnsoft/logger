package logger

import "io"

var stdLogger AbstractLogger = New()

func init() {
	stdLogger.SetReportCaller(true, DefaultSkipCallerNumber+1)
}

func GetDefault() AbstractLogger {
	return stdLogger
}

func SetDefault(logger AbstractLogger) {
	stdLogger = logger
}

func SetLogLevel(level Level) {
	stdLogger.SetLogLevel(level)
}

func SetOutput(writer io.Writer) {
	stdLogger.SetOutput(writer)
}

func SetDateFormat(format string) {
	stdLogger.SetDateFormat(format)
}

func SetReportCaller(b bool) {
	stdLogger.SetReportCaller(b, DefaultSkipCallerNumber+1)
}

func Debug(args ...interface{}) {
	stdLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	stdLogger.Debugf(format, args...)
}

func Print(args ...interface{}) {
	stdLogger.Print(args...)
}

func Printf(format string, args ...interface{}) {
	stdLogger.Printf(format, args...)
}

func Warning(args ...interface{}) {
	stdLogger.Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	stdLogger.Warningf(format, args...)
}

func Error(args ...interface{}) {
	stdLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	stdLogger.Errorf(format, args...)
}

func EntityLogger() AbstractLogger {
	return &LogEntity{Logger: stdLogger.(*Logger)}
}
