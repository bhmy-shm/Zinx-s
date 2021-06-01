package ziface

/*
 IRequest接口:
 把客户端请求的链接信息，和请求数据，包装到一个 Request中
*/

type IRequest interface {

	//获取当前连接
	GetConnection() IConnection

	//获取请求的消息数据
	GetData() []byte

	//获取请求的消息ID
	GetID() uint32
}
