# Zinx-S

Zinx-s 基于刘丹冰老师的 [Zinx](https://github.com/aceld/zinx) 二次开发，高效，可扩展，易于使用的go - Socket TCP网络通信框架。





## 功能介绍

- 提供多模块化布局
  - 快速建立TCP-Server句柄
  - 自定义路由业务
  - 消息封装与读写分离
  - Woker工作线程池
  - HOOK连接上下文钩子函数
- 添加功能
  - 公共组件（日志管理）
  - 插件化启动
  - 默认路由
  - 心跳监测，超时断开
  - 文件推送(protocol)
- 待实现功能
  - 断线重连
  - 自定义消息封装
  - 分布式
  - 过滤器



## 快速使用

**Server**

基于Zinx框架开发的服务器应用。完整的步骤如下：

  1. 创建需要加载服务端的功能
  2. 初始化Server句柄
  3. 注册路由处理业务
  4. 运行服务端

```go
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

	//6.注册路由
	//server.AddRouter(1,&HelloRouter{}) //注册第一个路由，支持注册多个
	server.DefaultRouter(&DefaultRouter{})	//默认路由不需要指定ID，ID由系统自动生成

	//7.运行服务端
	server.RunServer()
}
```



**其中自定义路由的方法如下：**

```go
type HelloRouter struct{
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest){
	fmt.Println("Call Router PreHandle")

	//先读取客户端的数据
	fmt.Println("recv from client: MsgID=",request.GetID(),"data=",string(request.GetData()))

	//更新心跳时间戳，不管是创建什么业务一定要加上这个心跳判断，可以封装成一个函数，添加过来
    //isUpdateTimeStamp(string(request.GetData()))
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

func isUpdateTimeStamp(res string){
 	if res == "" {
		fmt.Println("更新时间戳")
		request.GetConnection().UpdateTime()
	}   
}
```



---

**自定义日志记录**

在返回err错误的基础上，可以追加一些自己的语句到日志记录中。

```go
func(s *Server) Start(){
		//1.开启socket会话
		tcpSocket,err := net.ResolveTCPAddr("tcp","127.0.0.1:7777")
		errcode.PanicF(err,"ResolveTCPAddr is error")

		//2.开启Listener监听
		listen,err := net.ListenTCP("tcp",tcpSocket)
		errcode.PanicF(err,"ListenTCP is error")
	...
}

//Errcode 中定义的Panic恐慌记录日志
func PanicF(err error,format string) {
	if err != nil {
		global.Config.Logger.PanicF(err.Error(),format)
	}
}
```



---

**发送和接收通信数据**

默认的 '\n' 数据流边界收发方式，此方法读取时采用的是for{}循环单字节读取，按照\n 定义边界并 break

Server端方法如下：

```go
func main(){
    ...
    conn,err := listen.AcceptTCP()
	errcode.PanicF(err,"Listen Accept is error")
    
    //如果不使用封包的形式，就是默认的'\n' 边界方式
    bytes := utils.NewServerBytes(conn)

    for{	
	//读取对端发送数据
	str,err := bytes.ConnReadeBytes()
	if err != nil{
	    return
	}
        
    	//响应发送数据到对端
    	nts,err := bytes.ConnWriterBytes(s string)    
        if err != nil {
	    return
	}
    }
}
```

Client端方法如下：

```go
func main(){
    conn,err := net.Dial(c.IPVersion,"127.0.0.1:7777")
    errcode.PanicF(err,"Client Dial is Error")

    bytes := utils.NewClientBytes(conn)

	for{
            //向服务端发送请求，需要加上请求的客户端ID
	    _ , err := bytes.ConnWriteBytes(1,"hello world client")
	    if err != nil {
		return
	    }
            //读取服务端响应数据
	    str ,err := bytes.ConnReadeBytes()
	    if err != nil{
		return
	    }
	    fmt.Println("从服务端读取的数据：",str)
	}
}
```
## 配置文件

```yaml
Server:
  Name: Zinx-V.10-server
  Host: 127.0.0.1
  TcpPort: 7777
  MaxConn: 3
  HeartBeatTime: 5
  WorkerPoolSize: 10
```

`Name`:服务器应用名称

`Host`:服务器IP

`TcpPort`:服务器监听端口

`MaxConn`:允许的客户端链接最大数量

`HeartBeatTime:` 5

`WorkerPoolSize`:工作任务池最大工作Goroutine数量


