package errcode

import (
	"ZinX-shm/global"
	"io"
)

const(
	IOF string = "error = io.EOF,send is out"
)

//1.处理err返回值
func ErrorF(err error,format string){
	if err != nil {
		global.Config.Logger.Info(err.Error()+format)
	}
}

//2.处理io.EOF
func ErrorIOF(err error,format string) {
	if err == io.EOF {
		global.Config.Logger.WarnF(format)
	}else if err != io.EOF && err != nil {
		global.Config.Logger.ErrorF(format)
	}
}

//3.Panic恐慌
func PanicF(err error,format string) {
	if err != nil {
		global.Config.Logger.PanicF(format)
	}
}
