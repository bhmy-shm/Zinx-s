package ziface

import "net"

type IClient interface {
	//1.启动客户端
	Start()
	//2.停止客户端
	Stop()
	//3.运行客户端
	RunClient()
	//4.以'\n'的数据边界形式运行客户端
	DefaultClient()
	//
	SendHearBeatResponse(dp IDataPack,conn net.Conn)
}
