package znet

import (
	"ZinX-shm/utils"
	"ZinX-shm/utils/errcode"
	"ZinX-shm/ziface"
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct{

	//3.客户端的ID号
	ClientID uint32

	//1.连接的服务端ip版本
	IPVersion string

	//2.连接的服务端Host主机
	ServerHost string
}

func NewClient(id uint32,host string) *Client{
	return &Client{ClientID: id,IPVersion: "tcp",ServerHost: host}
}

func (c *Client) Start(){

	conn,err := net.Dial(c.IPVersion,"127.0.0.1:7777")
	errcode.PanicF(err,"Client Dial is Error")

	for{
		dp := NewDataPack()

		//封包写
		binaryMsg,err:=dp.Pack(NewMessagePackage(0,[]byte("ZinxV0.9 client 0 Test Message")))
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
			msg := msgHead.(*Message)
			msg.Data  = make([]byte,msg.GetMsgLen())

			if _,err := io.ReadFull(conn,msg.Data) ; err != nil{
				fmt.Println("read MsgData error",err)
				break
			}

			//响应心跳
			if string(msg.Data) == "ping" {
				c.SendHearBeatResponse(dp,conn)
			}

			fmt.Println("Recv Server Msg：ID=",msg.ID,"len=",msg.Len,"data=",string(msg.Data))
		}

		time.Sleep(3*time.Second)
	}
}

func (c *Client) Stop(){}

func (c *Client) RunClient(){
	c.Start()
}

func (c *Client) SendHearBeatResponse(dp ziface.IDataPack,conn net.Conn){
	binaryMsg,err:=dp.Pack(NewMessagePackage(101,[]byte("")))
	if err != nil{
		fmt.Println("Package error",err)
		return
	}

	_,_= conn.Write(binaryMsg)

}

func (c *Client) DefaultClient(){

	conn,err := net.Dial(c.IPVersion,"127.0.0.1:7777")
	errcode.PanicF(err,"Client Dial is Error")

	for{
		bytes := utils.NewClientBytes(conn)

		_,err := bytes.ConnWriteBytes(1,"hello world client")
		if err != nil {
			return
		}

		str ,err := bytes.ConnReadeBytes()
		if err != nil{
			return
		}
		fmt.Println("从服务端读取的数据：",str)


		time.Sleep(time.Second*3)
	}
}