package parsers

import (
	"IMWebServer/maps"
	"IMWebServer/models"
	"IMWebServer/models/IM_BaseDefine"
	"IMWebServer/models/IM_Buddy"
	"IMWebServer/models/IM_Message"
	"IMWebServer/models/IM_Server"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	// "sync"
)

type Message struct {
	Head Header          `json:"header"`
	Body json.RawMessage `json:"body"`
}

type Header struct {
	CommandId uint16 `json:"command_id"`
	ServiceId uint16 `json:"service_id"`
}

type Login struct {
	AccessToken   string `json:"access_token"`
	OnlineStatus  int    `json:"online_status"`
	ClientType    int    `json:"client_type"`
	DeviceId      string `json:"device_id"`
	ClientVersion string `json:"client_version"`
}

func Exec(conn *websocket.Conn) {
	for {
		m := &Message{}
		err := conn.ReadJSON(m)
		if err != nil {
			log.Println("read:", err)
			conn.WriteMessage(1, []byte(err.Error()))
			return
		}
		DataFilter(m, conn)
	}
	conn.Close()
}

//以下代码会在调试通过的时候进行重构
func DataFilter(m *Message, conn *websocket.Conn) (uint32, error) {
	ph := &PduHeader{Command_id: m.Head.CommandId, Service_id: m.Head.ServiceId}
	switch m.Head.CommandId {
	case 530:
		////获取最近联系人列表及未读消息记录数
		RecentContactList(ph, m.Body, conn)
	case 268:
		//登录
		LoginSys(ph, m.Body, conn)
	case 777:
		//获取聊天记录
		GetRecord(ph, m.Body, conn)
	case 769:
		//发送文本消息
		SendTextMsg(ph, m.Body, conn)
	}

	// if m.Head.CommandId == 530 {
	// 	//获取最近联系人列表及未读消息记录数
	// 	rcs := &IM_Buddy.IMRecentContactSessionReq{}
	// 	json.Unmarshal(m.Body, rcs)
	// 	ph := &PduHeader{Command_id: m.Head.CommandId, Service_id: m.Head.ServiceId}
	// 	d, _ := ph.RenderByte(rcs)
	// 	maps.Register(strconv.Itoa(int(*rcs.UserId)), conn)
	// 	models.DBChanSend <- d
	// } else if m.Head.CommandId == 268 {
	// 	//登录
	// 	l := &Login{}
	// 	json.Unmarshal(m.Body, l)
	// 	atv := &IM_Server.IMAccessTokenValidateReq{
	// 		AccessToken: proto.String(l.AccessToken),
	// 		AttachData:  []byte("1"),
	// 	}
	// 	ph := &PduHeader{Command_id: 0x0735, Service_id: m.Head.ServiceId}
	// 	d, _ := ph.RenderByte(atv)
	// 	maps.Register(*atv.AccessToken, conn)
	// 	models.DBChanSend <- d

	// } else if m.Head.CommandId == 777 {
	// 	//获取聊天记录
	// 	record := &Msgrecord{}
	// 	json.Unmarshal(m.Body, record)
	// 	var session_type IM_BaseDefine.SessionType = 1
	// 	msglist := &IM_Message.IMGetMsgListReq{
	// 		UserId:      proto.Uint32(record.UserId),
	// 		SessionType: &session_type,
	// 		SessionId:   proto.Uint32(record.SessionId),
	// 		MsgIdBegin:  proto.Uint32(record.MsgIdBegin),
	// 		MsgCnt:      proto.Uint32(record.MsgCnt),
	// 		AttachData:  []byte(record.AttachData),
	// 	}
	// 	ph := &PduHeader{Command_id: m.Head.CommandId, Service_id: m.Head.ServiceId}
	// 	d, _ := ph.RenderByte(msglist)
	// 	maps.Register(strconv.Itoa(int(*msglist.UserId)), conn)
	// 	models.DBChanSend <- d
	// } else if m.Head.CommandId == 769 {
	// 	//发送文本消息
	// 	msgdata := &IM_Message.IMMsgData{}
	// 	json.Unmarshal(m.Body, msgdata)
	// 	ph := &PduHeader{Command_id: m.Head.CommandId, Service_id: m.Head.ServiceId}
	// 	d, _ := ph.RenderByte(msgdata)
	// 	models.DBChanSend <- d
	// 	maps.Register(strconv.Itoa(int(*msgdata.FromUserId)), conn)

	// } else if m.Head.CommandId == 769 {
	// 	//发送图片消息

	// } else if strconv.FormatInt(int64(m.Head.CommandId), 16) == "0x0106" {
	// 	//退出登录

	// } else if strconv.FormatInt(int64(m.Head.CommandId), 16) == "0x0303" {
	// 	// 消息已读确认
	// 	mdr := &IM_Message.IMMsgDataReadAck{}
	// 	json.Unmarshal(m.Body, mdr)
	// 	ph := &PduHeader{Command_id: m.Head.CommandId, Service_id: m.Head.ServiceId}
	// 	d, _ := ph.RenderByte(mdr)
	// 	models.DBChanSend <- d
	// }

	return 0, nil
}

func RecentContactList(p *PduHeader, body []byte, conn *websocket.Conn) {
	//获取最近联系人列表及未读消息记录数
	rcs := &IM_Buddy.IMRecentContactSessionReq{}
	json.Unmarshal(body, rcs)
	d, _ := p.RenderByte(rcs)
	maps.Register(strconv.Itoa(int(*rcs.UserId)), conn)
	models.DBChanSend <- d
}

func LoginSys(p *PduHeader, body []byte, conn *websocket.Conn) {
	//登录
	l := &Login{}
	json.Unmarshal(body, l)
	atv := &IM_Server.IMAccessTokenValidateReq{
		AccessToken: proto.String(l.AccessToken),
		AttachData:  []byte("1"),
	}
	d, _ := p.RenderByte(atv)
	maps.Register(*atv.AccessToken, conn)
	models.DBChanSend <- d
}

func GetRecord(p *PduHeader, body []byte, conn *websocket.Conn) {
	//获取聊天记录
	record := &Msgrecord{}
	json.Unmarshal(body, record)
	var session_type IM_BaseDefine.SessionType = 1
	msglist := &IM_Message.IMGetMsgListReq{
		UserId:      proto.Uint32(record.UserId),
		SessionType: &session_type,
		SessionId:   proto.Uint32(record.SessionId),
		MsgIdBegin:  proto.Uint32(record.MsgIdBegin),
		MsgCnt:      proto.Uint32(record.MsgCnt),
		AttachData:  []byte(record.AttachData),
	}
	d, _ := p.RenderByte(msglist)
	maps.Register(strconv.Itoa(int(*msglist.UserId)), conn)
	models.DBChanSend <- d
}

func SendTextMsg(p *PduHeader, body []byte, conn *websocket.Conn) {
	//发送文本消息
	msgdata := &IM_Message.IMMsgData{}
	json.Unmarshal(body, msgdata)
	d, _ := p.RenderByte(msgdata)
	maps.Register(strconv.Itoa(int(*msgdata.FromUserId)), conn)
	models.DBChanSend <- d
}
