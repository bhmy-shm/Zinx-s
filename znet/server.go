package znet

import (
	"ZinX-shm/configs"
	"ZinX-shm/global"
	"ZinX-shm/utils"
	"ZinX-shm/utils/errcode"
	"ZinX-shm/ziface"
	"fmt"
	"net"
)

type Server struct {
	//3.服务端port
	Port int
	//1.服务端名称
	Name string
	//2.服务端ip
	IP string
	//4.服务端IP协议版本
	IPVersion string

	//5.启动插件接口
	conf *configs.Config

	//6.使用IMsgHandler替换Router
	Handler ziface.IMsgHandler

	//读Server创建连接之后自动调用的Hook函数 -- OnConnStart()
	OnConnStart func(conn ziface.IConnection)
	//读Server销毁连接之前自动调用的Hook函数 -- OnConStop()
	OnConnStop func(conn ziface.IConnection)
}

func NewServer(port int,options ...configs.Options) (*Server,error){

	server := &Server{
		Port: port,
		Name: global.GlobalSection.Name,
		IP: global.GlobalSection.Host,
		conf: global.Config,
		IPVersion: "ipv4",
		Handler: NewMsgQueueHandle(),
	}

	/* 初始化可插拔启动配置组件 */
	for _,opt := range options{
		if opt != nil {
			if err := opt.Apply(global.Config); err != nil{
				return nil,err
			}
		}
	}
	return server,nil
}

//1.启动服务器
func(s *Server) Start(){

	fmt.Printf("[Zinx] ServerName: %s Listen at IP:%s Port: %d is startning\n",global.GlobalSection.Name,
		global.GlobalSection.Host,global.GlobalSection.TcpPort)
	fmt.Printf("[ZinX] Version: %s, MaxConn:%d, MaxPackeetSize:%d \n",global.GlobalSection.Version,
		global.GlobalSection.MaxConn,global.GlobalSection.MaxPackageSize)

	go func() {
		//开启工作池
		s.Handler.StartWorkerPool()

		//1.开启socket会话
		tcpSocket,err := net.ResolveTCPAddr("tcp","127.0.0.1:7777")
		errcode.PanicF(err,"ResolveTCPAddr is error")

		//2.开启Listener监听
		listen,err := net.ListenTCP("tcp",tcpSocket)
		errcode.PanicF(err,"ListenTCP is error")


		var ConnID uint32 = 1
		for {
			//3.处理连接Socket会话的客户端请求
			conn,err := listen.AcceptTCP()
			errcode.PanicF(err,"Listen Accept is error")

			//4.将客户端的连接会话封装到连接模块中
			newConn := NewConnection(s,ConnID,conn,s.Handler)
			ConnID++

			//5.判断当前最大连接个数是否超出了最大连接数量，如果超过，则关闭这个新的连接
			if s.conf.Session.Len() >= global.GlobalSection.MaxConn {
				_ = conn.Close()
				continue
			}

			//6.将客户端连接加入到session中
			if s.conf.Session != nil {
				s.conf.Session.Add(ConnID,newConn)
			}else{
				ss := NewSession(5)
				ss.Add(ConnID,newConn)
			}

			//6.启动连接模块
			newConn.Start()
		}
	}()
}

//2.停止服务器
func(s *Server) Stop(){
	s.conf.Session.ClearConn()
}

//3.运行服务器
func(s *Server) RunServer(){
	s.Start()
	select {}
}

//4.添加Router
func(s *Server)	AddRouter(msgID uint32,router ziface.IRouter){
	s.Handler.AddRouter(msgID,router)
}

//5.添加默认Router
func(s *Server) DefaultRouter(router ziface.IRouter){
	msgID := uint32(utils.RandID())
	s.Handler.AddRouter(msgID,router)
}

/* ------------------------ 钩子函数 -------------------------- */
//注册OnConnStart() 钩子函数的方法
func(s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)){
	s.OnConnStart = hookFunc
}
//注册OnConnStop() 钩子函数的方法
func(s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)){
	s.OnConnStop = hookFunc
}

//调用OnConnStart() 钩子函数的方法
func(s *Server) CallOnConnStart(connection ziface.IConnection){
	if s.OnConnStart != nil {
		fmt.Println("------->Call OnConnStart()")
		s.OnConnStart(connection)
	}
}
//调用OnConnStop() 钩子函数的方法
func(s *Server) CallOnConnStop(connection ziface.IConnection){
	if s.OnConnStop != nil {
		fmt.Println("------->Call OnConnStop()")
		s.OnConnStop(connection)
	}
}