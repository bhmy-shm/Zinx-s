package znet

import (
	"ZinX-shm/global"
	"ZinX-shm/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

//封包，拆包的模块
type DataPack struct {}

//初始化方法
func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

//获取数据包头的长度
func (d *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节) + DataID uint32(4字节) = 8字节
	return 8
}


//封包方法
func(d *DataPack)Pack(msg ziface.IMessage)([]byte,error){

	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将数据长度写进databuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgLen()) ; err != nil{
		return nil,err
	}
	//将数据ID写进databuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgID()) ; err != nil{
		return nil,err
	}

	//将数据内容写入databuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgData()) ; err != nil{
		return nil,err
	}

	return dataBuff.Bytes(),nil

}


//拆包方法
func (d *DataPack) UnPack( binaryData []byte)(ziface.IMessage,error){
	//创建一个输入二进制数据的io.Reader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息，得到datalen和MsgID
	msg := &Message{}

	//读数据包头的长度
	if err := binary.Read(dataBuff,binary.LittleEndian,&msg.Len) ; err != nil{
		return nil,err
	}
	//赌数据包头的ID
	if err := binary.Read(dataBuff,binary.LittleEndian,&msg.ID) ; err != nil{
		return nil,err
	}
	//判断数据包的datalen是否已经超出了配置文件中设置的最大包长度
	if global.GlobalSection.MaxPackageSize > 0 && msg.Len > global.GlobalSection.MaxPackageSize{
		return nil,errors.New("too Large mswg data recv!")
	}

	return msg,nil
}
