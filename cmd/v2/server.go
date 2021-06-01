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

	//3.创建连接会话功能，加载心跳监测,传入心跳监测的超时时间
	heartBeat := znet.NewSession(5)

	//4.初始化Server端
	server,_ := znet.NewServer(7777,&configs.Config{
		Logger: newLogger,
		DataPack: dataPack,
		Session: heartBeat,
	})

	//5.调用Hook函数
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionAfter)

	//6.注册路由
	//server.AddRouter(1,&HelloRouter{}) //注册第一个路由，支持注册多个
	//server.AddRouter(2,&RBACRouter{}) //执行业务必须注册路由，并填写ID号
	server.DefaultRouter(&DefaultRouter{})	//默认路由不需要指定ID，ID由系统自动生成

	//7.运行服务端
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

type RBACRouter struct{
	znet.BaseRouter
}

func (this *RBACRouter) PreHandle(request ziface.IRequest){
	fmt.Println("Call RBACRouter PreHandle ")

	//先读取客户端的数据
	fmt.Println("recv from client: MsgID=",request.GetID(),"data=",string(request.GetData()))

	//更新心跳时间戳
	if string(request.GetData()) == "" {
		fmt.Println("更新时间戳")
		request.GetConnection().UpdateTime()
	}

	//先读取客户端的数据，再回写ping...Ping
	err := request.GetConnection().SendMsg(200,[]byte("rbac---shm is root"))
	if err != nil {
		fmt.Println(err)
	}
}

type DefaultRouter struct{
	znet.BaseRouter
	znet.Session
}

func (this *DefaultRouter) PreHandle(request ziface.IRequest){
	fmt.Println("Call DefaultRouter PreHandle ")

	//先读取客户端的数据
	fmt.Println("recv from client: MsgID=",request.GetID(),"data=",string(request.GetData()))

	//更新心跳时间戳
	if string(request.GetData()) == "" {
		fmt.Println("更新时间戳")
		request.GetConnection().UpdateTime()
	}

	//先读取客户端的数据，再回写ping...Ping
	err := request.GetConnection().SendMsg(200,[]byte("default router And Zinx"))
	if err != nil {
		fmt.Println(err)
	}

}

//创建一些钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("---> DoConnectionBegin is Called....")

	//在连接建立之前发送一些数据
	if err := conn.SendMsg(210,[]byte("DoConnectionBegin")) ; err != nil {
		fmt.Println(err)
	}

	//给当前的链接设置一些属性
	fmt.Println("Set Conn Property...")
	conn.SetProperty("Name","shm")
	conn.SetProperty("Github","http://baidu.com")
	conn.SetProperty("Home","www.bhmy.top/blog")
	conn.SetProperty("Blog","www.jianshu.com")
}

//创建连接断开前需要执行的钩子函数
func DoConnectionAfter(conn ziface.IConnection){
	//为其他用户广播下线
	fmt.Println("--->DoConnectionLost is Called...")
	fmt.Println("connID=",conn.GetConnID(),"is Lost....")

	if value ,err := conn.GetProperty("Name") ; err == nil {
		fmt.Println("Name=",value)
	}
	if value ,err := conn.GetProperty("Github") ; err == nil {
		fmt.Println("Github=",value)
	}
	if value ,err := conn.GetProperty("Home") ; err == nil {
		fmt.Println("Home=",value)
	}
	if value ,err := conn.GetProperty("Blog") ; err == nil {
		fmt.Println("Blog=",value)
	}
}