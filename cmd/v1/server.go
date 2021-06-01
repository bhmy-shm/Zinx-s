package main

import (
	"ZinX-shm/configs"
	"ZinX-shm/utils/logger"
	"ZinX-shm/ziface"
	"ZinX-shm/znet"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)


func main(){

	//启动服务器时加载功能插件(软功能)

	//1.加载日志模块
	newLogger := logger.New(&lumberjack.Logger{
		Filename: "storage/logs",	//记录日志的位置
		MaxSize:   600,  //设置日志文件允许的最大占用空间 600MB
		MaxAge:    10,   //日志文件的最大声明周期 10天
		LocalTime: true, //日志文件名的时间格式为本地时间
	},"",log.LstdFlags).WithCaller(2)	//找两层栈调用进行记录

	//2.加载封包拆包模块
	dataPack := znet.NewDataPack()

	//3.加载心跳监测,传入心跳监测的超时时间
	heartBeat := znet.NewSession(5)

	//3.初始化Server端
	server,_ := znet.NewServer(7777,&configs.Config{
		Logger: newLogger,
		DataPack: dataPack,
		Session: heartBeat,
	})

	server.AddRouter(0,&HelloRouter{}) //注册第一个路由，支持注册多个
	server.RunServer()
}

type HelloRouter struct{
	znet.BaseRouter
}

func (this *HelloRouter) PreHandle(request ziface.IRequest){
	fmt.Println("Call Router PreHandle")

	//先读取客户端的数据
	fmt.Println("recv from client: MsgID=",request.GetID(),"data=",string(request.GetData()))

	//更新心跳时间戳
	if string(request.GetData()) == "" {
		fmt.Println("更新时间戳")
		request.GetConnection().UpdateTime()
	}

	//先读取客户端的数据，再回写ping...Ping
	err := request.GetConnection().SendMsg(200,[]byte("hello....hi"))
	if err != nil {
		fmt.Println(err)
	}

}

//func (this *HelloRouter) Handle(request ziface.IRequest){
//	fmt.Println("Call Router Handle")
//	conn := request.GetConnection().GetTCPConnAmount()
//	bytes := utils.NewServerBytes(conn)
//	_,_ =bytes.ConnWriteBytes("ping... ping ...")
//}

//func (this *HelloRouter) PostHandle(request ziface.IRequest){
//	fmt.Println("Call Router PostHandle")
//	conn := request.GetConnection().GetTCPConnAmount()
//
//	bytes := utils.NewServerBytes(conn)
//	_,_ =bytes.ConnWriteBytes("after... ping ...")
//}

