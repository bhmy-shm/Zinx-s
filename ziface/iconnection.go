package ziface

import (
	"net"
)

/*
	客户端连接模块
	封装客户端连接状态和读取，发送数据的方法
*/

type IConnection interface {
	//1.启动连接
	Start()
	//2.关闭连接
	Stop()
	//3.获取当前的连接对象
	GetTCPConnAmount() *net.TCPConn
	//4.获取当前的连接对象的ID号
	GetConnID() uint32
	//5.获取当前的连接的ip+port
	GetRemoteAddr() net.Addr
	//6.读取当前连接数据
	ReadData()

	//获取时间戳
	TimeStamp() int64
	//更新连接的时间戳
	UpdateTime()

	//7.向客户端返回当前连接数据
	SendData()
	//8.封包写入数据
	SendMsg(msgID uint32, data []byte) error
	//9.拆包读取数据
	ReadMsg(conn *net.TCPConn) (IMessage,[]byte,error)

	/* ------- 设置Hook连接属性 ------ */
	//设置属性
	SetProperty(key string,value interface{})
	//获取属性
	GetProperty(key string)(interface{},error)
	//删除属性
	RemoveProperty(key string)
}


