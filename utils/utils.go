package utils

import (
	"ZinX-shm/global"
	"bytes"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type IBytes interface {
	ConnReadeBytes() (string,error)
	ConnWriteBytes(id uint32,s string)(int,error)
}

type ServerBytes struct{
	conn net.Conn
}
func NewServerBytes(conn net.Conn) *ServerBytes{
	return &ServerBytes{conn: conn}
}

type ClientBytes struct {
	conn net.Conn
}
func NewClientBytes(conn net.Conn) *ClientBytes{
	return &ClientBytes{conn: conn}
}


//按照字节读取,1个字节1个字节的读取
//边界：\n
func (this *ServerBytes) ConnReadeBytes() (string,error) {
	readBytes:= make([]byte,1)
	var buf bytes.Buffer
	//循环读取每次只读取1字节
	//如果读取的字节==\n 代表边界，跳出读取
	for{
		_,err := this.conn.Read(readBytes)
		if err != nil{
			return "",err
		}
		readByte := readBytes[0]
		if readByte == '\n' {
			break
		}
		buf.WriteByte(readByte)
	}
	return buf.String(),nil
}

//按照字节发送数据，定义边界：\n
func (this *ServerBytes) ConnWriteBytes(s string)(int,error){
	var buf bytes.Buffer

	//写入传入的数据
	buf.WriteString(s)
	buf.WriteByte('\n')

	//按照[]byte切片发送
	return this.conn.Write(buf.Bytes())
}


//按照字节读取,1个字节1个字节的读取
//边界：\n
func (this *ClientBytes) ConnReadeBytes() (string,error) {

	readBytes:= make([]byte,1)
	var buf bytes.Buffer
	//循环读取每次只读取1字节
	//如果读取的字节==\n 代表边界，跳出读取
	for{
		_,err := this.conn.Read(readBytes)
		if err != nil{
			return "",err
		}
		readByte := readBytes[0]
		if readByte == '\n' {
			break
		}
		buf.WriteByte(readByte)
	}
	return buf.String(),nil
}

//按照字节发送数据，定义边界：\n
func (this *ClientBytes) ConnWriteBytes(id uint32,s string)(int,error){
	var buf bytes.Buffer
	ids := strconv.Itoa(int(id))
	//写入传入的数据
	buf.WriteString(ids+s)
	buf.WriteByte('\n')

	//按照[]byte切片发送
	return this.conn.Write(buf.Bytes())
}


/*-------- 生成随机ID -------*/
func RandID() int {
	rand.Seed(time.Now().Unix())
	id := rand.Intn(global.GlobalSection.MaxConn)
	return id
}