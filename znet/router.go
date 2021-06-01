package znet

import "ZinX-shm/ziface"

type BaseRouter struct{}

//1.在处理Connection业务之前的方法Hook
func(b *BaseRouter) PreHandle(request ziface.IRequest){}

//2.再处理Connection的主方法
func(b *BaseRouter) Handle(request ziface.IRequest){}

//3.再处理Connection业务之后的方法Hook
func(b *BaseRouter) PostHandle(request ziface.IRequest){}

//4.创建一个专门处理返回心跳的路由
func(b *BaseRouter) ReadHeartBeat (request ziface.IRequest){}