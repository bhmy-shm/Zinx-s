package ziface

type IServer interface {
	//1.启动服务器
	Start()
	//2.停止服务器
	Stop()
	//3.运行服务器
	RunServer()
	//4.注册一个路由方法，为客户端连接提供业务处理
	AddRouter(msgID uint32,router IRouter)
	//5.注册默认路由
	DefaultRouter(router IRouter)

	/* --------- 钩子函数 -------- */
	//注册OnConnStart() 钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	//注册OnConnStop() 钩子函数的方法
	SetOnConnStop(func(connection IConnection))

	//调用OnConnStart() 钩子函数的方法
	CallOnConnStart(connection IConnection)
	//调用OnConnStop() 钩子函数的方法
	CallOnConnStop(connection IConnection)

}

