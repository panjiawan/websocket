package dto

type WebSocketSendReq struct {
	Platform string   `json:"platform"` // web, app
	ToUsers  []string `json:"to_users"` // 接收者 userid列表
	Message  string   `json:"message"`  // 消息内容
}
