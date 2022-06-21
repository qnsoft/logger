package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/fatih/color"
	"golang.org/x/term"
)

type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var defaultWriter = color.Output

const DefaultDateFormat = "2006/01/02 15:04"

const DefaultSkipCallerNumber = 3

var (
	risk  = []byte(":")
	left  = []byte("[")
	right = []byte("]")
	space = []byte(" ")
	brk   = []byte("\n")

	bufPool sync.Pool
)

type AbstractLogger interface {
	SetLogLevel(level Level)
	SetOutput(writer io.Writer)
	SetReportCaller(b bool, skipCallerNumber ...int)
	SetDateFormat(format string)

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	Print(args ...interface{})
	Printf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})

	EntityLogger() AbstractLogger
}

type Logger struct {
	io.Writer
	Level            Level
	DateFormat       string
	RecordCaller     bool
	SkipCallerNumber int

	entityData map[string]interface{}
}

func init() {
	bufPool = sync.Pool{New: func() interface{} {
		return bytes.NewBuffer(nil)
	}}
}

func New() *Logger {
	return &Logger{
		Writer:           defaultWriter,
		Level:            DebugLevel,
		DateFormat:       DefaultDateFormat,
		SkipCallerNumber: DefaultSkipCallerNumber,
	}
}

func (l *Logger) SetOutput(writer io.Writer) {
	l.Writer = writer
}

func (l *Logger) SetLogLevel(level Level) {
	l.Level = level
}

func (l *Logger) SetDateFormat(format string) {
	l.DateFormat = format
}

func (l *Logger) getArgs(args []interface{}) []interface{} {
	argWithSpaces := make([]interface{}, len(args)*2)
	idx := 0
	for _, v := range args {
		argWithSpaces[idx] = v
		idx++
		argWithSpaces[idx] = " "
		idx++
	}
	return argWithSpaces
}

func (l *Logger) SetReportCaller(b bool, skipCallerNumber ...int) {
	l.RecordCaller = b
	if len(skipCallerNumber) == 0 {
		l.SkipCallerNumber = DefaultSkipCallerNumber
	} else if skipCallerNumber[0] > 0 {
		l.SkipCallerNumber = skipCallerNumber[0]
	}
}

func (l *Logger) Debug(args ...interface{}) {
	if l.Level <= DebugLevel {
		l.WriteString(DebugLevel, fmt.Sprint(l.getArgs(args)...))
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Level <= DebugLevel {
		l.WriteString(DebugLevel, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Print(args ...interface{}) {
	if l.Level <= InfoLevel {
		l.WriteString(InfoLevel, fmt.Sprint(l.getArgs(args)...))
	}
}

func (l *Logger) Printf(format string, args ...interface{}) {
	if l.Level <= InfoLevel {
		l.WriteString(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Warning(args ...interface{}) {
	if l.Level <= WarnLevel {
		l.WriteString(WarnLevel, fmt.Sprint(l.getArgs(args)...))
	}
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	if l.Level <= WarnLevel {
		l.WriteString(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Error(args ...interface{}) {
	stack := debug.Stack()
	args = append(args, "\n", *(*string)(unsafe.Pointer(&stack)))
	l.WriteString(ErrorLevel, fmt.Sprint(l.getArgs(args)...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	format += "\n %s"
	args = append(args, string(debug.Stack()))
	l.WriteString(ErrorLevel, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Error(args...)
	panic(fmt.Sprint(args...))
}

func (l *Logger) WriteString(level Level, message string) {
	bytesBuf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		bytesBuf.Reset()
		bufPool.Put(bytesBuf)
	}()
	if len(l.DateFormat) > 0 {
		bytesBuf.Write(left)
		bytesBuf.WriteString(time.Now().Format(l.DateFormat))
		bytesBuf.Write(right)
		bytesBuf.Write(space)
	}
	if !l.isTerminal() {
		bytesBuf.WriteString(defaultFormatters[level].Type)
	} else {
		bytesBuf.WriteString(defaultFormatters[level].ColorType)
	}
	bytesBuf.Write(space)
	l.writeCallerInfo(bytesBuf)
	bytesBuf.WriteString(message)
	if l.entityData != nil {
		byts, _ := json.Marshal(&l.entityData)
		bytesBuf.Write(space)
		bytesBuf.Write(byts)
	}
	bytesBuf.Write(brk)

	_, _ = l.Writer.Write(bytesBuf.Bytes())
}

func (l *Logger) writeCallerInfo(buf *bytes.Buffer) {
	if l.RecordCaller {
		_, callerFile, line, ok := runtime.Caller(l.SkipCallerNumber)
		if ok {
			buf.Write([]byte(path.Base(callerFile)))
			buf.Write(risk)
			buf.WriteString(strconv.Itoa(line))
			buf.Write(risk)
			buf.Write(space)
		}
	}
}

func (l *Logger) isTerminal() bool {
	switch v := l.Writer.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func (l *Logger) EntityLogger() AbstractLogger {
	entity := &LogEntity{Logger: l}
	entity.LoggerIdField = "TraceLoggerId"
	return entity
}
