package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}
type RestResponse struct {
	Result []Update `json:"result"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]string `json:"keyboard"`
	ResizeKeyboard  bool       `json:"resize_keyboard"`
	OneTimeKeyboard bool       `json:"one_time_keyboard"`
}

type BotMessage struct {
	ChatId      int                  `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup *ReplyKeyboardMarkup `json:"reply_markup,omitempty"`
}
