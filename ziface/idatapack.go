package ziface

/*
	封包，拆包，模块；直接面向TCP链接中的数据流，用于处理TCP粘包问题
*/

type IDataPack interface {
	//获取包头的长度方法
	GetHeadLen() uint32
	//封包
	Pack(msg IMessage)([]byte,error)
	//拆包
	UnPack([]byte)(IMessage,error)
}

