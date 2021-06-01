package ziface

/*
	消息管理抽象层，对路由层进行封装
*/

type IMsgHandler interface {

	//1.调度执行路由器，调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)

	//2.添加路由器，为消息添加具体的处理业务逻辑
	AddRouter(msgID uint32,router IRouter)

	//3.启动worker工作池
	StartWorkerPool()

	//4.将消息发送给任务队列进行处理
	SendMsgToTaskQueue(request IRequest)
}