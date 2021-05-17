package main

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/smtp"
)

func SendEmail(message, to string) {
	// Sender data.
	from := "redioteka@internet.ru"
	from_from := "CoolRedTech"

	// Receiver email address.
	sendTo := []string{
		to,
	}

	// smtp server configuration.
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"

	// Message.

	// Authentication.
	auth := smtp.PlainAuth("", from, from_from, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, sendTo, []byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func FormFilmList(movieTitles []string) string {
	res := ""
	for _, film := range movieTitles {
		res += film + "\n"
	}
	return res
}

// domain.user with only id and email
func GetSendInfo(db *database.DBManager) ([]domain.User, error) {
	return nil, nil
}

func GetWatchlist(userId uint, db *database.DBManager) ([]string, error) {
	return nil, nil
}

func SendWatchList(movieTitles []string, sendTo string) error {
	message := FormFilmList(movieTitles)
	SendEmail(message, sendTo)
	return nil
}

func SendWatchlists(db *database.DBManager) error {
	users, err := GetSendInfo(db)
	if err != nil {
		return err
	}

	for _, user := range users {
		watchlist, err := GetWatchlist(user.ID, db)
		if err != nil {
			log.Err(err)
		}
		err = SendWatchList(watchlist, user.Email)
		if err != nil {
			log.Err(err)
		}
	}
	return nil
}

func main() {
	db := database.Connect()
	err := SendWatchlists(db)
	if err != nil {
		log.Err(err)
	}
}
