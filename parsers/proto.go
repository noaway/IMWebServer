package parsers

import (
	// "IMWebServer/constant"
	"IMWebServer/maps"
	"IMWebServer/models"
	"IMWebServer/models/IM_BaseDefine"
	"IMWebServer/models/IM_Buddy"
	"IMWebServer/models/IM_Login"
	"IMWebServer/models/IM_Message"
	"IMWebServer/models/IM_Server"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/golang/protobuf/proto"
	"log"
	"strconv"
	"time"
	// "errors"
)

func ProtoInit() {
	go func() {
		for data := range models.DBChanRece {
			r := &PduHeader{}
			h, d, err := r.GetPduHeader(data)
			if h == nil {
				continue
			}
			//一下代码会在调试通过的时候进行重构
			if err != nil {
				log.Println(err.Error())
			}
			if h.Command_id != 1793 {
				log.Println(h.Command_id)
				// log.Println(strconv.FormatInt(int64(0x030a), 10))
			}

			if h.Command_id == 1793 {
				// log.Println(h)
			} else if h.Command_id == 531 {
				//获取最近联系人列表及未读消息记录数
				rensession := &IM_Buddy.IMRecentSesAndUnreadCntRsp{}
				proto.Unmarshal(d, rensession)
				log.Println(rensession.ContactSessionList)

				msg := &Message{}
				msg.Head.CommandId = h.Command_id
				msg.Head.ServiceId = h.Service_id
				msg.Body, _ = json.Marshal(rensession)

				if conn := maps.Conns(strconv.Itoa(int(*rensession.UserId))); conn != nil {
					conn.WriteJSON(msg)
				}

			} else if h.Command_id == 1796 {
				//登录
				ValidateRsp := &IM_Server.IMValidateRsp{}
				if err := proto.Unmarshal(d, ValidateRsp); err != nil {
					log.Println("错误是:", err)
					// continue
				}
				// log.Println(ValidateRsp)
				// if *ValidateRsp.ResultCode == 0 {
				// 	if conn := Maps.Conns(*ValidateRsp.AccessToken); conn != nil {
				// 		conn.WriteJSON(ValidateRsp)
				// 	}
				// 	return
				// }
				var ResultCode IM_BaseDefine.ResultType = 0
				var UserStatType IM_BaseDefine.UserStatType = 0

				LoginRes := &IM_Login.IMLoginRes{
					ServerTime:   proto.Uint32(uint32(time.Now().Unix())),
					ResultCode:   &ResultCode,
					ResultString: proto.String(""),
					OnlineStatus: &UserStatType,
					UserInfo:     ValidateRsp.UserInfo,
					AccessToken:  ValidateRsp.AccessToken,
				}

				msg := &Message{}
				msg.Head.CommandId = h.Command_id
				msg.Head.ServiceId = h.Service_id
				msg.Body, _ = json.Marshal(LoginRes)

				if conn := maps.Conns(*LoginRes.AccessToken); conn != nil {
					conn.WriteJSON(msg)
				}

			} else if h.Command_id == 778 {
				//获取聊天记录
				msglist := &IM_Message.IMGetMsgListRsp{}
				if err := proto.Unmarshal(d, msglist); err != nil {
					log.Println(err)
					continue
				}
				log.Println(msglist.MsgList)
				if conn := maps.Conns(strconv.Itoa(int(*msglist.UserId))); conn != nil {
					conn.WriteJSON(msglist)
				}
			} else if h.Command_id == 769 {
				//发送文本消息
				msgdata := &IM_Message.IMMsgData{}
				if err := proto.Unmarshal(d, msgdata); err != nil {
					log.Println(err)
					continue
				}
				sessiontype := IM_BaseDefine.SessionType_SESSION_TYPE_SINGLE

				ack := &IM_Message.IMMsgDataAck{
					UserId:      msgdata.FromUserId,
					SessionId:   msgdata.ToSessionId,
					MsgId:       msgdata.MsgId,
					SessionType: &sessiontype,
				}

				d, _ := h.RenderByte(msgdata)
				models.RouteChan <- d

				if conn := maps.Conns(strconv.Itoa(int(*ack.UserId))); conn != nil {
					conn.WriteJSON(ack)
				}
			}
		}
	}()

	// go func() {
	// 	for data := range models.RouteChan {
	// 		r := &PduHeader{}
	// 		h, _, _ := r.GetPduHeader(data)
	// 		log.Println("route send: ", h)
	// 	}
	// }()
}

type PduHeader struct {
	Length     uint32 `json:"length"`
	Version    uint16 `json:"version"`
	Flag       uint16 `json:"flag"`
	Service_id uint16 `json:"service_id"`
	Command_id uint16 `json:"command_id"`
	Seq_num    uint16 `json:"seq_num"`
	Reversed   uint16 `json:"reversed"`
}

func (self *PduHeader) GetPduHeader(data []byte) (*PduHeader, []byte, error) {
	defer func() {
		if err := recover(); err != nil {
			beego.Error(err)
		}
	}()
	head := data[:16]
	bytesBuffer := bytes.NewBuffer(head)
	if err := binary.Read(bytesBuffer, binary.BigEndian, self); err != nil {
		return nil, nil, err
	}
	head = []byte("")
	return self, data[16:], nil
}

func (self *PduHeader) RenderByte(protobuf proto.Message) ([]byte, error) {
	if protobuf == nil {
		bf := new(bytes.Buffer)
		self.Length = uint32(16)
		if err := binary.Write(bf, binary.BigEndian, self); err != nil {
			log.Println(err)
			return nil, err
		}
		return bf.Bytes(), nil
	}
	if data, err := proto.Marshal(protobuf); err != nil {
		log.Println(err)
		return nil, err
	} else {
		bf := new(bytes.Buffer)
		self.Length = uint32(len(data) + 16)
		if err := binary.Write(bf, binary.BigEndian, self); err != nil {
			log.Println(err)
			return nil, err
		}
		by := bf.Bytes()
		by = append(by, data...)
		return by, nil
	}
}

func (self *PduHeader) HeartbeatPacket() ([]byte, error) {
	self.Service_id = 0x0007
	self.Command_id = 0x0701
	return self.RenderByte(nil)
}
