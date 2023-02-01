package main

import (
	"log"
	"os"
	"tg-bot-p2p/pkg/api"
	"tg-bot-p2p/pkg/bot"
	"tg-bot-p2p/pkg/repository"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	token := os.Getenv("BOT_TOKEN")
	channelID := os.Getenv("BOT_CHANNEL")
	mongoConnString := os.Getenv("MONGO_CONN")

	tg := api.NewTelegramAPI(token, channelID)

	mongodb, err := repository.NewMongoClient(mongoConnString)

	if err != nil {
		log.Fatalln("Unable connect to database cause:", err.Error())
	}

	userSessionStorage := repository.NewUserSessionStorage(mongodb)
	dayStorage := repository.NewDayStorage(mongodb)
	botInstance := bot.NewBot(userSessionStorage, tg, dayStorage)

	updatesChan := tg.Poller.ListenUpdates()
	botInstance.StartHandling(updatesChan)

}
