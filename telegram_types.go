package telegoat

// https://core.telegram.org/bots/api#update

type GetUpdatesResponse struct {
	Result []Update
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	From      User   `json:"from"`
	Text      string `json:"text"`
}

type User struct {
	Id       int
	Username string
}
