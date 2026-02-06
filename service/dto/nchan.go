package dto

type NChanEventDto struct {
	EventType string `json:"event_type"` // subscribe, unsubscribe, message
	ChannelID string `json:"channel_id"` // user_id
	Platform  string `json:"platform"`   // web app
}
