package ziface

/*
路由的抽象接口，路由里的封装的请求数据是基于IRequest
*/

type IRouter interface{
	//1.在处理Connection业务之前的方法Hook
	PreHandle(request IRequest)

	//2.再处理Connection的主方法
	Handle(request IRequest)

	//3.再处理Connection业务之后的方法Hook
	PostHandle(request IRequest)

	//4.处理心跳路由
	ReadHeartBeat(request IRequest)

}

