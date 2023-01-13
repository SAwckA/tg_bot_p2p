package api

import "strings"

type UpdateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type MessageResponse struct {
	OK       bool `json:"ok"`
	*Message `json:"result"`
}

type ErrorResponse struct {
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type Chat struct {
	ID int `json:"id"`

	// "private" / "group" / "supergroup" / "channel"
	Type string `json:"type"`
}

type User struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
}

type Message struct {
	MessageID int  `json:"message_id"`
	From      User `json:"from"`
	*Chat     `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

func (m *Message) IsCommand() bool {
	return strings.HasPrefix(m.Text, "/")
}

func (m *Message) Command() string {
	return strings.TrimPrefix(m.Text, "/")
}

type CallbackQuery struct {
	ID              string `json:"id"`
	From            User   `json:"from"`
	*Message        `json:"message"`
	InlineMessageID string `json:"inline_message_id"`
	ChatInstance    string `json:"chat_instance"`
	Data            string `json:"data"`
}

type Update struct {
	UpdateID       int `json:"update_id"`
	*Message       `json:"message"`
	EditedMessage  *Message `json:"edited_message"`
	ChannelPost    *Message `json:"channel_post"`
	*CallbackQuery `json:"callback_query"`
}
