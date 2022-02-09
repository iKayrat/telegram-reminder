package main

import (
	"log"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	conf "github.com/iKayrat/telegram-reminder/handlers"
)

func main() {
	config, err := conf.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env", err)
	}

	bot, err := tgbot.NewBotAPI(config.TokenApi)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)
			log.Printf("** [%s] %s", update.Message.Chat.FirstName, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
