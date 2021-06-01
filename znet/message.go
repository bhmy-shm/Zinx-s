package znet

type Message struct{
	//消息的ID
	ID uint32
	//消息的长度
	Len uint32
	//消息内容
	Data []byte
}

func NewMessagePackage(id uint32,data []byte) *Message{
	return &Message{
		ID: id,
		Len: uint32(len(data)),
		Data: data,
	}
}

//获取消息ID
func(m *Message)GetMsgID() uint32 {
	return m.ID
}
//获取消息的长度
func(m *Message)GetMsgLen() uint32 {
	return m.Len
}
//获取消息的内容
func(m *Message)GetMsgData() []byte {
	return m.Data
}

//设置消息ID
func(m *Message)SetMsgId(i uint32){
	m.ID = i
}
//设置消息的长度
func(m *Message)SetMsgLen(i uint32){
	m.Len = i
}
//设置消息的内容
func(m *Message)SetMsgData(data []byte){
	m.Data = data
}