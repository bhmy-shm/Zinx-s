package ziface

/*
	会话列表，用来推送，拉取客户端连接的数据
*/

type ISession interface {

	//添加连接
	Add(id uint32,connection IConnection)
	//删除连接
	Remove(connID uint32)
	//根据Conn-ID获取连接
	GetSession(connID uint32) (IConnection,error)
	//得到当前连接总数
	Len()int
	//清除并终止所有连接
	ClearConn()



	/* 心跳监测 */
	//1.发送心跳
	Send(id uint32)
	//2.心跳监测判断
	HeartBeat(id uint32,num int64)
	//3.向客户端推送服务端的文件数据
	Upload(file string)
}