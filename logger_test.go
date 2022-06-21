package logger

import (
	"io/ioutil"
	"testing"
)

func TestLogger(t *testing.T) {
	Error("This is an error message", "this is error message 2")
	Errorf("This is an error message %s", "这是一个错误信息")
	Warning("This is a warning message")
	Print("This is an info message")
	Debug("This is a debug message")
	entity := EntityLogger()
	entity.(*LogEntity).SetDataKV("name", "xiusin")
	entity.(*LogEntity).SetDataKV("age", 1)
	entity.(*LogEntity).SetDataKV("address", []string{"HeNan", "ZhengZhou", "JinShuiQu"})
	entity.Print("print")
}

func BenchmarkLogger(b *testing.B) {
	SetLogLevel(DebugLevel)
	f, _ := ioutil.TempFile("/Users/xiusin/projects/src/github.com/qnsoft/logger/", "*.log")
	SetOutput(f)
	defer f.Close()
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Debug("[%d] This is a debug message")
	}
}
