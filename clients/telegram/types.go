package telegram

type UpdateResponse struct {
	Ok bool `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID int `json:"update_id"`
	// поле Message может отсутствовать поэтому указываем тип как ссылку на структуру
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From `json:"from"`
	Chat Chat `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}