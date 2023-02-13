package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	_ "github.com/lib/pq"
)

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

	TOKEN := os.Getenv("TELEGRAM_TOKEN")
	DBSOURCE := os.Getenv("DBSOURCE")

	db := DBconnection(DBSOURCE)
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}

	bot, err := tgbot.NewBotAPI(TOKEN)
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
		err := db.QueryRow("SELECT days FROM workdays WHERE user_id = $1", userID).Scan(&days)
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

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}

func DBconnection(dbsource string) *sql.DB {
	// connStr := "user=pg dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", dbsource)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

// type Config struct {
// 	token          string `mapstructure:"TELEGRAM_TOKEN"`
// 	satabaseSource string `mapstructure:"DBSOURCE"`
// }

// func LoadConfig(path string) (config Config, err error) {
// 	viper.AddConfigPath(path)
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("env")

// 	viper.AutomaticEnv()

// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		return
// 	}
// 	err = viper.Unmarshal(&config)
// 	return
// }
