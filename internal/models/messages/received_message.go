package messages

type ReceivedMessage struct {
	Id       string `json:"id"`
	Body     string `json:"body"`
	Username string `json:"username"`
	Receiver string `json:"receiver"`
	IsNew    bool   `json:"is_new"`
}
