package controllers

import (
	"database/sql"
	"fmt"
	"log"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bio struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Type      string `json:"type"`
	Location  string `json:"location"`
}

func Start(db *sql.DB, info *tgbot.Chat) error {
	log.Printf("********info: %#v \n", info)
	log.Fatalf("********info: %#v \n", info)
	fmt.Printf("********infoBIO: %#v \n", info.Bio)

	bio := Bio{
		Firstname: info.FirstName,
		Lastname:  info.LastName,
		Username:  info.UserName,
		Type:      info.Type,
		Location:  info.Location.Address,
	}
	log.Fatal(bio)

	err := db.QueryRow(`INSERT INTO users (firstname,lastname,username,type,location) VALUES 
	($1,$2,$3,$4,$5)`, info.FirstName, info.LastName, info.UserName, info.Type, info.Location.Address).Scan()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
