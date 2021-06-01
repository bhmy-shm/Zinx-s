package main

import "ZinX-shm/znet"

func main(){

	//1.创建一个客户端
	client := znet.NewClient(1,"127.0.0.1:7777")

	//2.创建一个客户端路由向服务端发送数据，将发送的数据设置成一个结构体进行封装，然后一并发送，选择发送的

	client.RunClient()
}