package main

import (
	"tg-bot-p2p/pkg/api"
	"tg-bot-p2p/pkg/bot"
	"tg-bot-p2p/pkg/repository"
)

func main() {
	token := "1224606820:AAFiYvv8-eiMdwD7M_ypDbfTFjI3QyQW6pI"

	tg := api.NewTelegramAPI(token)

	userSessionStorage := repository.NewUserSessionStorage()
	botInstance := bot.NewBot(userSessionStorage, tg)

	updatesChan := tg.Poller.ListenUpdates()
	botInstance.StartHandling(updatesChan)

}
