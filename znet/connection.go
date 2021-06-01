package znet

import (
	"ZinX-shm/global"
	"ZinX-shm/utils"
	"ZinX-shm/utils/errcode"
	"ZinX-shm/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Connection struct{

	//当前Conneciton隶属于哪个Server
	TcpServer ziface.IServer

	//当前连接的状态
	IsClosed bool

	//Channel捕获退出的管道
	ExitChan chan bool

	//当前连接的客户端ID
	ConnID uint32

	//当前连接的客户端时间戳
	Timestamp int64

	//当前连接的socket套接字
	Conn *net.TCPConn

	//消息的管理MsgID，和对应的处理业务API关系
	Handler ziface.IMsgHandler

	//读写Goroutine之间消息通信的管道
	msgChan chan []byte

	//连接属性的集合
	Property map[string]interface{}

	//保护连接属性的锁
	PropertyLock sync.Mutex
}

func NewConnection(server ziface.IServer,id uint32,conn *net.TCPConn,handler ziface.IMsgHandler) *Connection {
	connection :=  &Connection{
		TcpServer: server,
		IsClosed: false,
		ConnID: id,
		Conn: conn,
		Timestamp: time.Now().Unix(),
		ExitChan: make(chan bool,1),
		msgChan: make(chan []byte),
		Property: make(map[string]interface{}),
		Handler: handler,
	}

	return connection
}

//返回当前时间戳的时间
func (c *Connection) TimeStamp() int64{
	return c.Timestamp
}

//更新时间函数
func (c *Connection) UpdateTime(){
	c.Timestamp = time.Now().Unix()
}

//1.启动连接
func(c *Connection) Start(){
	fmt.Println("Conn Start() ...ConnID",c.ConnID)

	//TODO 启动从当前连接 "读数据" 的业务goroutine
	go c.ReadData()

	//TODO 启动从当前连接 ”写数据“ 的业务goroutine
	go c.SendData()

	//TODO 按照开发者传递进来的，创建连接之后需要调用的处理业务，执行Hook函数
	c.TcpServer.CallOnConnStart(c)
}

//2.关闭连接
func(c *Connection) Stop(){
	fmt.Println("this Connection is Closed...")

	//关闭状态
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	//在销毁连接之前，需要调用关闭的Hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭Scoket连接
	_ = c.Conn.Close()


	//告诉ExitChan,ReadData关闭
	c.ExitChan <- true

	close(c.ExitChan)
	close(c.msgChan)
}

//3.获取当前的连接对象
func(c *Connection) GetTCPConnAmount() *net.TCPConn{
	return c.Conn
}
//4.获取当前的连接对象的ID号
func(c *Connection) GetConnID() uint32{
	return c.ConnID
}
//5.获取当前的连接的ip+port
func(c *Connection) GetRemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}
//6.读取当前连接数据
func(c *Connection) ReadData(){
	fmt.Println("Reader Goroutine is Running...")
	defer c.Stop()
	defer fmt.Println("ConnID =",c.ConnID,"Reader is exit,remote add is",c.GetRemoteAddr().String())


	//先判断是否使用封包的形式
	if global.Config.DataPack != nil {
		for {
			//拆包读取消息
			msg, data, err := c.ReadMsg(c.GetTCPConnAmount())
			if len(data) == 0 && err != nil {
				return
			}

			//然后记录到Message消息模块
			msg.SetMsgData(data)

			//再将得到消息，记录到Request请求模块
			req := Request{
				conn: c,
				msg:  msg,
			}

			//判断是否开启了工作池机制，将消息发送给Worker工作池处理
			if global.GlobalSection.WorkerPoolSize > 0 {
				go c.Handler.SendMsgToTaskQueue(&req)
			}else{
				//调用MsgHandler来处理业务，根据消息请求，到达MsgHandler，
				//寻找指定MsgID的router路由进行调用。
				go c.Handler.DoMsgHandler(&req)
			}
		}
	} else {
		for{
			//如果不使用封包的形式，就是默认的'\n' 边界方式
			bytes := utils.NewServerBytes(c.GetTCPConnAmount())

			//读取客户端请求数据
			str,err := bytes.ConnReadeBytes()
			if str == ""{
				return
			}
			if err != nil {
				errcode.ErrorF(err,"[connection.go 128]")
				return
			}
			fmt.Println("从客户端读到的请求：",str)

			//将请求数据封装到请求模块
			msg := NewMessagePackage(c.ConnID,[]byte(str))

			//再将得到消息，记录到Request请求模块
			req := Request{
				conn: c,
				msg: msg,
			}
			if global.GlobalSection.WorkerPoolSize > 0 {
				go c.Handler.SendMsgToTaskQueue(&req)
			}else{
				//调用MsgHandler来处理业务，根据消息请求，到达MsgHandler，
				//寻找指定MsgID的router路由进行调用。
				go c.Handler.DoMsgHandler(&req)
			}
		}
	}
}

//7.向客户端返回当前连接数据,写消息的goroutine发送给客户端消息
func(c *Connection) SendData(){
	fmt.Println("Writer Goroutine is Runnning...")
	defer fmt.Println(c.GetRemoteAddr().String(),"conn Writer exit!")

	//1.监控管道，循环阻塞的从管道中读取数据
	for{
		select{
		case data := <-c.msgChan:
			//2.如果有数据写给客户端
			if _,err := c.Conn.Write(data) ; err != nil {
				fmt.Println("Send data error:",err)
				return
			}
		//如果c.ExitChan 管道中可以读到数据了，代表Read已经退出，Write也要退出。
		case <- c.ExitChan :
			return
		}
	}

}

//封包写，写数据，提供一个SendMsg方法，将服务端要发送给客户端的数据，先进行封包，再发送
func(c *Connection) SendMsg(msgID uint32,data []byte) error {

	if c.IsClosed == true {
		global.Config.Session.Remove(msgID)	//如果连接关闭了，就要调用全局删除这条会话记录
		return errors.New("Connection is Closed when send msg")
	}

	//将Data进行封包 MsgDataLen | MsgDataID | MsgData
	dp := NewDataPack()
	msg := NewMessagePackage(msgID,data)

	//拿到二进制准备发送的数据
	binaryMsg ,err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Pack error msg id = ",msgID)
		return errors.New("Package error msg")
	}

	//将拿到的数据发送给msgChan
	c.msgChan <- binaryMsg

	return nil
}

//拆包读，
func (c *Connection) ReadMsg(conn *net.TCPConn) (ziface.IMessage,[]byte,error){

	headData := make([]byte,global.Config.DataPack.GetHeadLen())
	if _,err := io.ReadFull(conn,headData) ; err == io.EOF || err != nil {
		return nil,nil,err
	}

	//拆包得到MsgID，MsgDataLen，放到一个message消息对象中
	msg,err := global.Config.DataPack.UnPack(headData)
	if err != nil{
		fmt.Println("unpack error",err)
		return nil,nil,err
	}

	//根据DataLen,再次读取Data，放到message.Data这个字段中
	var data []byte
	if msg.GetMsgLen() > 0 {
		data = make([]byte,msg.GetMsgLen())
		if _,err := io.ReadFull(conn,data) ; err != nil {
			fmt.Println("read msg data error",err)
			return nil,nil,err
		}
	}
	//fmt.Println("data=",string(data))
	return msg,data,nil
}


/* --------------------  设置连接属性 ------------------- */
//设置连接属性
func (c *Connection)SetProperty(key string,value interface{}){
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()

	c.Property[key] = value
}

//获取连接属性
func (c *Connection)GetProperty(key string)(interface{},error){
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()

	if value,found := c.Property[key] ; found {
		return value,nil
	}else{
		return nil,errors.New("No Property found")
	}
}

//移除连接属性
func (c *Connection)RemoveProperty(key string){
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()

	delete(c.Property,key)
}

/*
dp := NewDataPack()
	//先判断是否使用封包的形式
	if dp != nil {
		for {
			//拆包读取消息
			msg, data, err := c.ReadMsg(c.GetTCPConnAmount(),dp)
			if len(data) == 0 && err != nil {
				return
			}

			//然后记录到Message消息模块
			msg.SetMsgData(data)

			//再将得到消息，记录到Request请求模块
			req := Request{
				conn: c,
				msg:  msg,
			}

			//调用路由来处理业务，从路由中找到注册绑定的Conn对应的router调用
			go func(request ziface.IRequest) {
				c.Router.PostHandle(request)
				c.Router.Handle(request)
				c.Router.PreHandle(request)
				c.Router.ReadHeartBeat(request)
			}(&req)
		}
	} else {
		for{
			//如果不使用封包的形式，就是默认的'\n' 边界方式
			bytes := utils.NewServerBytes(c.GetTCPConnAmount())

			//读取客户端请求数据
			str,err := bytes.ConnReadeBytes()
			if str == ""{
				return
			}
			if err != nil {
				errcode.ErrorF(err,"[connection.go 128]")
				return
			}
			fmt.Println("从客户端读到的请求：",str)

			//将请求数据封装到请求模块
			msg := NewMessagePackage(c.ConnID,[]byte(str))

			//再将得到消息，记录到Request请求模块
			req := Request{
				conn: c,
				msg: msg,
			}
			//调用路由来处理业务，从路由中找到注册绑定的Conn对应的router调用
			go func(request ziface.IRequest) {
				c.Router.PostHandle(request)
				c.Router.Handle(request)
				c.Router.PreHandle(request)
				c.Router.ReadHeartBeat(request)
			}(&req)
		}
	}

*/

/*
	for {
		//创建一个拆包解包的对象
		dp := NewDataPack()

		//读取客户端的 MsgHead 二进制流8个字节
		headData := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(c.GetTCPConnAmount(),headData) ; err == io.EOF || err != nil {
			break
		}

		//拆包得到MsgID，MsgDataLen，放到一个message消息对象中
		msg,err := dp.UnPack(headData)
		if err != nil{
			fmt.Println("unpack error",err)
			break
		}

		//根据DataLen,再次读取Data，放到message.Data这个字段中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte,msg.GetMsgLen())
			if _,err := io.ReadFull(c.GetTCPConnAmount(),data) ; err != nil {
				fmt.Println("read msg data error",err)
				break
			}
		}
		msg.SetMsgData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg: msg,
		}

		go func(request ziface.IRequest) {
			c.Router.PostHandle(request)
			c.Router.Handle(request)
			c.Router.PreHandle(request)
			c.Router.ReadHeartBeat(request)
		}(&req)
	}*/