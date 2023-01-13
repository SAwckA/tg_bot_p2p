package main

import (
	"github.com/joho/godotenv"
	"os"
	"tg-bot-p2p/pkg/api"
	"tg-bot-p2p/pkg/bot"
	"tg-bot-p2p/pkg/repository"
)

func main() {

	_ = godotenv.Load()

	token := os.Getenv("BOT_TOKEN")

	tg := api.NewTelegramAPI(token)

	userSessionStorage := repository.NewUserSessionStorage()
	botInstance := bot.NewBot(userSessionStorage, tg)

	updatesChan := tg.Poller.ListenUpdates()
	botInstance.StartHandling(updatesChan)

}
