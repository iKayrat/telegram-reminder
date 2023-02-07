package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iKayrat/telegram-reminder/controllers"

	_ "github.com/lib/pq"
)

const chatID = 118469838

var numericKeyboard = tgbot.NewReplyKeyboard(
	tgbot.NewKeyboardButtonRow(
		tgbot.NewKeyboardButton("1"),
		tgbot.NewKeyboardButton("2"),
		tgbot.NewKeyboardButton("3"),
	),
	tgbot.NewKeyboardButtonRow(
		tgbot.NewKeyboardButton("4"),
		tgbot.NewKeyboardButton("5"),
		tgbot.NewKeyboardButton("6"),
	),
)

func main() {
	ctx := context.Background()

	conf, err := controllers.LoadConfig(".")
	if err != nil {
		log.Panic("cannot load env", err)
	}
	log.Println(conf)

	db := controllers.DBconnection(conf.DatabaseSource)
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	bot, err := tgbot.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		// Get user ID and message text
		userID := update.Message.From.ID
		text := update.Message.Text

		// Check if user is already registered
		var days int
		// err := db.QueryRow("SELECT days FROM workdays WHERE user_id = ?", userID).Scan(&days)
		err := db.QueryRowContext(ctx, "SELECT days FROM workdays WHERE user_id = $1", userID).Scan(&days)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			log.Fatal(err)
			continue
		}

		// Register user if they have not registered yet
		if err == sql.ErrNoRows {
			_, err := db.Exec("INSERT INTO workdays (user_id, days) VALUES ($1, 0)", userID)
			if err != nil {
				log.Println(err)
				continue
			}
			days = 0
		}

		msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch text {
		case "/start":
			msg := tgbot.NewMessage(update.Message.Chat.ID, "Welcome! Use /work to track your workdays.")
			bot.Send(msg)

		case "/work":
			days++
			_, err := db.Exec("UPDATE workdays SET days = $1 WHERE user_id = $2", days, userID)
			if err != nil {
				log.Println(err)
				continue
			}

			msg := tgbot.NewMessage(update.Message.Chat.ID, "Great!")
			bot.Send(msg)

		case "/status":
			day := 0
			err := db.QueryRow("SELECT days FROM workdays WHERE user_id = $1", userID).Scan(&day)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				continue
			}
			// str := "You have worked " + string(day) + " days"
			str := fmt.Sprintf("You have worked %d days", day)
			msg := tgbot.NewMessage(update.Message.Chat.ID, str)
			bot.Send(msg)

		case "/open":
			msg.ReplyMarkup = numericKeyboard
		case "/close":
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		}

		// if _, err := bot.Send(msg); err != nil {
		// 	log.Panic(err)
		// }
	}
}
