package models

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type CreateMessageReq struct {
	RoomID  string `json:"room_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type GetListMessagesReq struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

type GetListMessagesRes struct {
	Messages []Message
}
