package znet

import (
	"ZinX-shm/services"
	"ZinX-shm/ziface"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

//---------------------------------------------------SESSION管理类------------------------------------------------------

type Session struct {
	//每一个连接的会话
	sessions map[uint32]ziface.IConnection
	//拿到心跳的间隔时间
	num      int64
	//读写锁
	lock     sync.Mutex
}

func NewSession(num int64) *Session {
	return &Session{sessions: make(map[uint32]ziface.IConnection,0),num: num}
}

//添加会话
func (this *Session) Add(id uint32,connection ziface.IConnection) {
	this.lock.Lock()
	defer this.lock.Unlock()

	//向会话集合中追加一个当前连接的集合
	if _,found := this.sessions[id] ; !found {
		this.sessions[id] = connection
	}
	fmt.Println("connection Add to Session successful: conn num=",this.Len())

	this.runHeartBeat(id)
}

func(this *Session) runHeartBeat(id uint32){
	//发送心跳
	go this.Send(id)

	//开启go程，判断是响应时间是否超时
	go this.HeartBeat(id,this.num)
}

//删除连接
func(this *Session) Remove(connID uint32){
	//加锁
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.IsExitsConn() {
		delete(this.sessions,connID)
		fmt.Println("connection remove to Session successful: conn num=",this.Len())
	}
	return
}

//根据Conn-ID获取连接
func(this *Session )GetSession(connID uint32) (ziface.IConnection,error){
	this.lock.Lock()
	defer this.lock.Unlock()

	if conn,ok := this.sessions[connID] ; ok {
		return conn,nil
	}else{
		return nil,errors.New("connection not Found!")
	}
}

//得到当前连接总数
func(this *Session) Len()int{
	return len(this.sessions)
}

//清除并终止所有连接
func(this *Session) ClearConn(){
	this.lock.Lock()
	defer this.lock.Unlock()

	//删除conn并停止conn的工作
	for connID,conn := range this.sessions{
		//停止
		conn.Stop()
		//删除
		if this.IsExitsConn() {
			delete(this.sessions,connID)
		}
	}
	fmt.Println("connection Clear to Session successful: conn num=",this.Len())
}

//判断当前会话中是否还有连接的conneciton
func (this *Session) IsExitsConn() bool {
	if len(this.sessions)  == 0 {
		return false
	}else{
		return true
	}
}


/* ----------------------- 心跳功能 ---------------------- */

//发送心跳
func (this *Session) Send(id uint32){
	defer fmt.Println("已经终止向该客户端发送心跳: ",id)
	//不断地向客户端发送心跳包，间隔5秒发送一次
	for {
		//封包发送心跳内容
		fmt.Println("----x----",id)
		err := this.sessions[id].SendMsg(id,[]byte("ping"))
		if err != nil {
			return
		}

		//发送之后并记录时间戳
		this.sessions[id].UpdateTime()

		//间隔5秒重复循环发送
		time.Sleep( time.Duration(this.num) * time.Second )
	}
}

//心跳检测判断  每秒遍历一次 查看所有sess 上次接收消息时间  如果超过 num 就删除该 sess
func (this *Session) HeartBeat(id uint32,num int64) {

	for {
		time.Sleep(time.Second)
		for _,v:= range this.sessions{
			//fmt.Println("start:",time.Now().Unix())
			//fmt.Println("stop",v.TimeStamp())
			if time.Now().Unix() - v.TimeStamp() > num {
				break
			}
		}
	}
}



/* -------------------- 推送服务端数据 -------------------- */
func(this *Session) Upload(file string) {

	//传入一个ID，传入一个文件名

	//将文件的内容读取出来，然后转换成[]byte记录在缓存中
	str := "测试protobuf文件推送"
	strs := []byte(str)

	//拿到这个记录的长度，内容，和ID后对Mes-proto进行赋值
	msg := &services.MsgRequest{
		Message_ID:701,
		Message_Len: int32(len(strs)),
		Message_Data: strs,
	}

	//编码：
	//将定义好的proto对象，进行序列化拿到二进制文件
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("marshal err:", err)
	}

	//向每一个Session会话中的练级发送二进制数据
	for _,conn := range this.sessions {
		err := conn.SendMsg(conn.GetConnID(),data)
		if err != nil {
			fmt.Println("推送数据出错")
		}
	}
}