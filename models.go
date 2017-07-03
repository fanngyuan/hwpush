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
	ActionType int               `json:"type"`
	Param      map[string]string `json:"param"`
}

type Message struct {
	MessageType int           `json:"type"`
	Body        MessageBody   `json:"body"`
	Action      MessageAction `json:"action"`
}

type Hps struct {
	Msg Message `json:"msg"`
	Ext Ext     `json:"ext"`
}

type Ext struct {
	Customize *[]Element `json:"customize"`
}

type Element map[string]string

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
		Param:      make(map[string]string),
	}
	message := Message{
		MessageType: 3,
		Body:        body,
		Action:      action,
	}
	customize := make([]Element, 0, 3)
	ext := Ext{
		Customize: &customize,
	}
	hps := Hps{
		Msg: message,
		Ext: ext,
	}
	return Notification{
		Hps: hps,
	}
}

func (this Notification) AddActionParam(key, value string) {
	this.Hps.Msg.Action.Param[key] = value
}

func (this Notification) AddExtra(key, value string) {
	valueMap := make(map[string]string)
	valueMap[key] = value
	eliment := Element(valueMap)
	*this.Hps.Ext.Customize = append(*this.Hps.Ext.Customize, eliment)
}

type Result struct {
	RequestId string `json:"requestId"`
	Msg       string `json:"msg"`
	Code      string `json:"code"`
}
