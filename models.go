package hwpush

type OauthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Error struct {
	ErrorCode   int    `json:"error"`
	Description string `json:"error_description"`
}

func (this Error) Error() string {
	return this.Description
}

type MessageBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type MessageAction struct {
	ActionType int `json:"type"`
	//Param      map[string]string `json:"param"`
}

type Message struct {
	MessageType int           `json:"type"`
	Body        MessageBody   `json:"body"`
	Action      MessageAction `json:"action"`
}

type Hps struct {
	Msg Message `json:"msg"`
	//Ext map[string]string `json:"ext"`
}

type Notification struct {
	Hps Hps `json:"hps"`
}

func NewNotification(content, title string) Notification {
	body := MessageBody{
		Title:   title,
		Content: content,
	}
	action := MessageAction{
		ActionType: 3,
	}
	message := Message{
		MessageType: 1,
		Body:        body,
		Action:      action,
	}
	hps := Hps{
		Msg: message,
	}
	return Notification{
		Hps: hps,
	}
}

type Result struct {
	RequestId string `json:"requestId"`
	Msg       string `json:"msg"`
	Code      string `json:"code"`
}
