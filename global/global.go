package global

import (
	"ZinX-shm/configs"
	"ZinX-shm/utils/Viper"
	"ZinX-shm/ziface"
)

var(
	Config = &configs.Config{}
	GlobalSection *section
)

//提供一个init方法，初始化当前的全局对象,如果全局配置文件没有加载，默认值
func init(){
	GlobalSection = &section{
		Name: "Zinx-Server",
		Version: "V0.9",
		TcpPort: 7777,
		Host: "0.0.0.0",
		MaxConn: 1000,
		MaxPackageSize: 4096,

		//一共有多少个负责处理worker链接池中业务的goroutine
		WorkerPoolSize: 10,

		//每个worker对应的消息队列任务数量的最大值，每个worker-channel能记录多少条消息记录，channle的容量
		MaxWorkerTaskLen: 1024,
	}
	//加载配置文件的开关
	GlobalSection.Reload()
}

//全局配置文件映射结构体
type section struct{
	//服务器Zinx框架数据包的最大值
	MaxPackageSize uint32

	//worker工作池的goroutine数量
	WorkerPoolSize uint32
	//ZinX框架允许用户开启的worker数量
	MaxWorkerTaskLen uint32

	//心跳检测的间隔时间
	HeartBeatTime int64

	//服务器主机监听端口号
	TcpPort int

	//服务器主机允许链接的最大值
	MaxConn int

	//当前Zinx全局的Server对象
	TcpServer ziface.IServer

	//服务器主机监听地址IP
	Host string
	//服务器名称
	Name string
	//服务器版本号
	Version string
}

//加载配置文件开关方法
func (s *section) Reload(){

	setting,err := Viper.NewSetting()
	if err != nil {
		panic(err)
	}
	err = setting.ReadSection("Server",s)
	if err != nil {
		panic(err)
	}
}
