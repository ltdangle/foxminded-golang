package bot

import (
	"time"

	"github.com/google/uuid"
)

// Telegram api response.
type TelegramResponse struct {
	Ok     bool     `json:"ok,omitempty"`
	Result []Update `json:"result"`
}

type Update struct {
	Message struct {
		Chat struct {
			FirstName string `json:"first_name,omitempty"`
			ID        int    `json:"id,omitempty"`
			Type      string `json:"type,omitempty"`
			Username  string `json:"username,omitempty"`
		} `json:"chat,omitempty"`
		Date     float64 `json:"date,omitempty"`
		Entities *[]struct {
			Length float64 `json:"length,omitempty"`
			Offset float64 `json:"offset,omitempty"`
			Type   string  `json:"type,omitempty"`
		} `json:"entities,omitempty"`
		Location *struct {
			Latitude  float64 `json:"latitude,omitempty"`
			Longitude float64 `json:"longitude,omitempty"`
		} `json:"location,omitempty"`
		From struct {
			FirstName    string  `json:"first_name,omitempty"`
			ID           float64 `json:"id,omitempty"`
			IsBot        bool    `json:"is_bot,omitempty"`
			LanguageCode string  `json:"language_code,omitempty"`
			Username     string  `json:"username,omitempty"`
		} `json:"from,omitempty"`
		MessageID float64 `json:"message_id,omitempty"`
		Text      *string `json:"text,omitempty"`
	} `json:"message,omitempty"`
	UpdateID int `json:"update_id,omitempty"`
}

// Bot message.
type BotMessage struct {
	ChatId      int                  `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup *ReplyKeyboardMarkup `json:"reply_markup,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]string `json:"keyboard"`
	ResizeKeyboard  bool       `json:"resize_keyboard"`
	OneTimeKeyboard bool       `json:"one_time_keyboard"`
}

func NewBotMessage(text string) *BotMessage {
	return &BotMessage{Text: text}
}

// Subscription.
type Subscription struct {
	Id        string    `json:"id"`
	ChatId    int       `json:"chat_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	FirsName  string    `json:"firstName"`
	Username  string    `json:"username"`
	Schedule  string    `json:"schedule"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSubscription() *Subscription {
	return &Subscription{Id: uuid.New().String(), UpdatedAt: time.Now()}
}
