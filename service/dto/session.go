package dto

type IsOnlineUserReq struct {
	Platform string   `json:"platform"`
	UserIds  []string `json:"user_ids"`
}

type GetOnlineUserCountReq struct {
	Platform string `form:"platform"`
}
