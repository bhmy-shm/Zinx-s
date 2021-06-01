package main

import (
	"ZinX-shm/ziface"
	"ZinX-shm/znet"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var wg sync.Mutex

/*模拟客户端*/
func main(){

	fmt.Println("Client Start....") ; time.Sleep(time.Second*1)

	//1.直接连接远程服务器，得到一个conn连接
	conn,err := net.Dial("tcp","127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err exit !")
	}

	//2.调用write方法写数据
	for {
		//发送封包的message数据
		dp := znet.NewDataPack()
		binaryMsg,err:=dp.Pack(znet.NewMessagePackage(2,[]byte("ZinxV0.9 client 1 Test Message")))
		if err != nil{
			fmt.Println("Package error",err)
			return
		}

		_,_= conn.Write(binaryMsg)

		//拆掉服务端响应的包

		//1.先读取流中的 Msghead,得到ID和datalen
		headData := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(conn,headData) ; err != nil {
			fmt.Println("read head error",err)
			break
		}

		//2.再根据DataLen进行第二次读取，将数据头拿出来
		msgHead,err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("Un pack error",err)
			break
		}

		//3.再根据数据头的长度进行读取，将data数据读取出来
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data  = make([]byte,msg.GetMsgLen())

			if _,err := io.ReadFull(conn,msg.Data) ; err != nil{
				fmt.Println("read MsgData error",err)
				break
			}

			//响应心跳
			if string(msg.Data) == "ping" {
				SendHearBeatResponse2(dp,conn)
			}
			fmt.Println("Recv Server Msg：ID=",msg.ID,"len=",msg.Len,"data=",string(msg.Data))
		}

		time.Sleep(3*time.Second)
	}
}


func SendHearBeatResponse2(dp ziface.IDataPack,conn net.Conn){
	wg.Lock()
	defer wg.Unlock()

	binaryMsg,err:=dp.Pack(znet.NewMessagePackage(101,[]byte("")))
	if err != nil{
		fmt.Println("Package error",err)
		return
	}
	_,_= conn.Write(binaryMsg)
}

