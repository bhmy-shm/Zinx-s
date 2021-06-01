package znet

import (
	"ZinX-shm/global"
	"ZinX-shm/ziface"
	"fmt"
)

type MsgHandler struct{
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter

	//负责worker取任务的消息队列,
	TaskQueue []chan ziface.IRequest

	//负责worker工作池的连接数量
	WorkerPoolSize uint32

	//默认随机ID
	defaultID uint32
}

//提供一个初始化创建MsgHandler方法
func NewMsgHandle() *MsgHandler{
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//创建消息队列的方法
func NewMsgQueueHandle() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
		WorkerPoolSize: global.GlobalSection.WorkerPoolSize,
		TaskQueue: make([]chan ziface.IRequest,global.GlobalSection.WorkerPoolSize),
	}
}


//1.调度执行路由器，调度/执行对应的Router消息处理方法
func(m *MsgHandler) DoMsgHandler(request ziface.IRequest) {

	//从request寻找map中指定request.message.ConnID的，看看有没有这个注册的路由方法
	handler,ok := m.Apis[m.defaultID]
	if !ok {
		//如果找不到代表这个MsgID还没有被注册路由
		fmt.Println("api MsgID = ",request.GetID(),"is Not Found ! Need Register router...")
	}

	fmt.Println("request:",request.GetID(),string(request.GetData()))

	//2.根据MsgID调用router对应的业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}


//2.添加路由器，为消息添加具体的处理业务逻辑
func(m *MsgHandler)AddRouter(msgID uint32,router ziface.IRouter) {

	//如果当前ID被注册了，就不需要添加Router了
	if _,found := m.Apis[msgID] ; !found {
		m.Apis[msgID] = router
	}

	m.defaultID = msgID

	fmt.Println("Add api MsgID=",msgID,"is successful !")
}

//启动一个worker工作池(开启工作池的动作一个框架只能发生一次开启工作池的动作)
func(m *MsgHandler) StartWorkerPool(){

	//拿到配置文件中指定的Worker数量，进行循环，然后依次启动一个worker
	for i:=0; i< int(m.WorkerPoolSize); i++{
		//1.启动每一个worker对应的channel作为载体,并初始化开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest,global.GlobalSection.MaxWorkerTaskLen)

		//2.启动当前的Worker，阻塞等待消息数据产生，从channel中传递进来
		go m.startOneWorker(i,m.TaskQueue[i] )
	}

}

//启动一个worker工作流程（）
func(m *MsgHandler) startOneWorker(workerID int, tastQueue chan ziface.IRequest){
	fmt.Println("WorkerID=",workerID,"is started")
	//不断阻塞等待对应的消息队列消息
	for{
		select {
		//如果有消息过来，就是一个客户端的消息请求，将其传给DoMsg专门处理Request中携带的任务。
		case req := <-tastQueue:
			m.DoMsgHandler(req)
		}
	}
}

//将消息发送给消息队列，由Worker进行处理
func(m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest){

	//1.将消息平均分配给不同的Worker
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("AddConnID=",request.GetConnection().GetConnID(),"request MsgID=",request.GetID(),"to workerID=",workerID)

	//2.将消息发送给对应的worker的TaskQueue
	m.TaskQueue[workerID] <- request
}