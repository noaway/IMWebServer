package parsers

type Msgrecord struct {
	//777
	UserId      uint32 `json:"user_id"`
	SessionType int32  `json:"session_type"`
	SessionId   uint32 `json:"session_id"`
	MsgIdBegin  uint32 `json:"msg_id_begin"`
	MsgCnt      uint32 `json:"msg_cnt"`
	AttachData  string `json:"attach_data"`
}
