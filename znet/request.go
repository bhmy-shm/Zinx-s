package znet

import "ZinX-shm/ziface"


type Request struct{
	//已经和客户端建立好的连接Connection
	conn ziface.IConnection
	//客户端请求的数据
	msg ziface.IMessage
}

//获取当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

//获取请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

//获取请求的消息ID
func (r *Request) GetID() uint32 {
	return r.conn.GetConnID()
}