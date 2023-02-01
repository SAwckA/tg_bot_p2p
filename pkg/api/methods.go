package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	sendMessage   = "sendMessage"
	editMessage   = "editMessageText"
	deleteMessage = "deleteMessage"
)

type Button struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type Keyboard [][]Button

type InlineKeyboard struct {
	Keyboard Keyboard `json:"inline_keyboard"`
}

type sendMessageReq struct {
	ChatID      int            `json:"chat_id"`
	Text        string         `json:"text"`
	ReplyMarkup InlineKeyboard `json:"reply_markup"`
}

type sendMessageReqStringID struct {
	ChatID      string         `json:"chat_id"`
	Text        string         `json:"text"`
	ReplyMarkup InlineKeyboard `json:"reply_markup"`
}

type editMessageReq struct {
	ChatID      int            `json:"chat_id"`
	MessageID   int            `json:"message_id"`
	Text        string         `json:"text"`
	ReplyMarkup InlineKeyboard `json:"reply_markup"`
}

func (tg *TelegramAPI) SendToChannel(text string, keyboard Keyboard) *Message {

	values := sendMessageReqStringID{
		ChatID:      tg.channelID,
		Text:        text,
		ReplyMarkup: InlineKeyboard{keyboard},
	}

	jsonData, err := json.Marshal(values)

	if err != nil {
		return nil
	}

	url := fmt.Sprintf("%s", urlAPIMethod(tg.token, sendMessage))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	body, err := ioutil.ReadAll(resp.Body)

	var data MessageResponse
	var errorData ErrorResponse

	err = json.Unmarshal(body, &data)

	err = json.Unmarshal(body, &errorData)

	fmt.Println(errorData)

	return data.Message
}

func (tg *TelegramAPI) SendMessage(chatID int, text string, keyboard Keyboard) *Message {

	values := sendMessageReq{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: InlineKeyboard{keyboard},
	}

	jsonData, err := json.Marshal(values)

	if err != nil {
		return nil
	}

	url := fmt.Sprintf("%s", urlAPIMethod(tg.token, sendMessage))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	body, err := ioutil.ReadAll(resp.Body)

	var data MessageResponse
	var errorData ErrorResponse

	err = json.Unmarshal(body, &data)

	err = json.Unmarshal(body, &errorData)

	fmt.Println(errorData)

	return data.Message
}

func (tg *TelegramAPI) EditMessage(chatID int, messageID int, text string, keyboard Keyboard) *Message {

	values := editMessageReq{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        text,
		ReplyMarkup: InlineKeyboard{keyboard},
	}

	jsonData, err := json.Marshal(values)

	if err != nil {
		return nil
	}

	url := urlAPIMethod(tg.token, editMessage)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	body, err := ioutil.ReadAll(resp.Body)

	var data MessageResponse
	var errorData ErrorResponse
	err = json.Unmarshal(body, &data)

	if data.OK == false {
		err = json.Unmarshal(body, &errorData)
	}

	return data.Message
}

func (tg *TelegramAPI) DeleteMessage(chatID int, messageID int) {
	url := urlAPIMethod(tg.token, deleteMessage)

	body := map[string]int{"chat_id": chatID, "message_id": messageID}
	jsonBody, err := json.Marshal(body)

	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		fmt.Println("Error during delete message")
	}

}
