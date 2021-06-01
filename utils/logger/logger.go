package logger

import (
	"fmt"
	"io"
	"log"
	"runtime"
)

type Level int8

type Fields map[string]interface{}	//记录日志信息的map表

const(
	LevelDebug Level = iota	//全部
	LevelInfo	//常规
	LevelWarn	//警告
	LevelError	//错误
	LevelFatal	//严重错误
	LevelPanic	//恐慌
)

//作为字符串的形式返回
func(l Level) String() string{
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

type LogInterface interface {
	clone() *Logger
	WithLevel(lvl Level) *Logger
	WithFields(f Fields) *Logger
	WithCaller(skip int) *Logger
	WithCallersFrames()  *Logger
	JSONFormat(message string) map[string]interface{}
	Output(message string)
	Debug( v ...interface{})
	DebugF( format string, v ...interface{})
	Info( v ...interface{})
	InfoF( format string, v ...interface{})
	Warn( v ...interface{})
	WarnF( format string, v ...interface{})
	Error( v ...interface{})
	ErrorF( format string, v ...interface{})
	Fatal( v ...interface{})
	FatalF( format string, v ...interface{})
	Panic( v ...interface{})
	PanicF( format string, v ...interface{})
}

//日志标准化结构体
type Logger struct{
	newLogger *log.Logger	//日志对象，i/o操作
	level Level				//日志等级
	fields Fields			//记录日志的map表
	callers []string		//记录栈的文件名，函数，行号等信息
}

//初始化调用日志对象
func New(w io.Writer,prefix string,flag int)  *Logger {
	l := log.New(w,prefix,flag)	//系统级日志对象调用
	return &Logger{newLogger: l}
}

//临时日志对象
func (l *Logger) clone() *Logger{
	nl := l
	return nl
}

//设置日志等级
func (l *Logger) WithLevel(lvl Level) *Logger{
	ll := l.clone()
	ll.level = lvl
	return ll
}

//设置日志公共字段
func(l *Logger) WithFields(f Fields) *Logger{
	ll := l.clone()
	if ll.fields == nil{
		ll.fields = make(Fields)	//初始化map
	}
	for k,v := range f{		//将传入的日志信息写入map
		ll.fields[k] = v
	}
	return ll
}

//记录调用日志的文件名，行号，执行函数
func(l *Logger) WithCaller(skip int)*Logger{
	ll := l.clone()
	pc,file,line,ok := runtime.Caller(skip)	//报告goroutine栈上调用的文件名和行号
	if ok {
		f := runtime.FuncForPC(pc)	//报告栈上所使用的函数
		ll.callers = []string{fmt.Sprintf("%s:%d %s",file,line,f.Name())}
	}
	return ll
}
//
func(l *Logger) WithCallersFrames() *Logger{
	maxCallerDepth := 25
	minCallerDepth := 1

	callers := []string{}
	pcs := make([]uintptr,maxCallerDepth)

	//返回函数/文件/行信息
	depth := runtime.Callers(minCallerDepth,pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for frame,more  := frames.Next();more;frame,more = frames.Next(){
		callers = append(callers,fmt.Sprintf("%s:%d %s",frame.File,frame.Line,frame.Function))
		if !more { break }
	}

	ll := l.clone()
	ll.callers = callers
	return ll
}