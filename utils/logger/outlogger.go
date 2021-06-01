package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

/* --------------------------------- 日志格式化 ------------------------------------- */

//将所有要输出的日志信息按照json格式化组合在一起
func (l *Logger) JSONFormat(message string) map[string]interface{}{
	data := make(Fields,len(l.fields)+4)

	data["level"] = l.level.String()
	data["time"] = time.Now().Local().String()
	data["message"] = message
	data["callers"] = l.callers

	if len(l.fields) > 0 {
		for k,v := range l.fields {
			if _,ok := data[k] ; !ok {
				data[k] = v
			}
		}
	}
	return data
}

//然后根据不同的日志级别，进行输出
func(l *Logger) Output(message string) {
	//输出前要将组合的日志信息map,转变成字符串
	body,_ := json.Marshal(l.JSONFormat(message))
	content := string(body)

	//然后根据不同的日志级别，选择不同的log.Print / log.Fatal 输出方式
	switch l.level {
	case LevelDebug:
		l.newLogger.Print(content)
	case LevelInfo:
		l.newLogger.Print(content)
	case LevelWarn:
		l.newLogger.Print(content)
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Fatal(content)
	case LevelPanic:
		l.newLogger.Panic(content)
	}
}

/* -------------------------------- 日志分级别输出 ------------------------------------- */

//根据先前定义的日志分级，编写对应的日志输出的外部方法，继续写入如下代码：
func (l *Logger) Debug(v ...interface{}) {
	l.WithLevel(LevelDebug).Output(fmt.Sprint(v...))
}

func (l *Logger) DebugF(format string, v ...interface{}) {
	l.WithLevel(LevelDebug).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Info( v ...interface{}) {
	l.WithLevel(LevelInfo).Output(fmt.Sprint(v...))
}

func (l *Logger) InfoF(format string, v ...interface{}) {
	l.WithLevel(LevelInfo).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.WithLevel(LevelWarn).Output(fmt.Sprint(v...))
}

func (l *Logger) WarnF(format string, v ...interface{}) {
	l.WithLevel(LevelWarn).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.WithLevel(LevelError).Output(fmt.Sprint(v...))
}

func (l *Logger) ErrorF(format string, v ...interface{}) {
	l.WithLevel(LevelError).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.WithLevel(LevelFatal).Output(fmt.Sprint(v...))
}

func (l *Logger) FatalF(format string, v ...interface{}) {
	l.WithLevel(LevelFatal).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(v ...interface{}) {
	l.WithLevel(LevelPanic).Output(fmt.Sprint(v...))
}

func (l *Logger) PanicF(format string, v ...interface{}) {
	l.WithLevel(LevelPanic).Output(fmt.Sprintf(format, v...))
}
