package main

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/cast"
	"fmt"
	"log"
	"net/smtp"

	"gopkg.in/gomail.v2"
)

func SendEmailGomail(message, to string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "Redioteka <redioteka@internet.ru>")
	m.SetHeader("")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Вы хотели посмотреть эти штуки!")
	m.SetBody("text/plain", message)

	d := gomail.NewDialer("smtp.mail.ru", 587, "redioteka@internet.ru", "CoolRedTech")

	// Send the email to Bob, Cora and Dank.
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func SendEmailSmtp(message, to string) {
	// Sender data.
	from := "redioteka@internet.ru"
	fromPass := "CoolRedTech"

	// Receiver email address.
	sendTo := []string{
		to,
	}

	// smtp server configuration.
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, fromPass, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, sendTo, []byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func FormFilmList(movieTitles []domain.Movie) string {
	res := ""
	for _, movie := range movieTitles {
		res += movie.Title + "\n"
	}
	return res
}

const querySelectWatchlist = "select id, email from users;"

func GetSendInfo(db *database.DBManager) ([]domain.User, error) {
	data, err := db.Query(querySelectWatchlist)
	if err != nil {
		return nil, err
	}
	var users []domain.User
	for _, row := range data {
		users = append(users, domain.User{
			ID:    cast.ToUint(row[0]),
			Email: cast.ToString(row[1]),
		})
	}
	return users, nil
}

const querySelectUserWatchlist = `select distinct m.id, m.title, m.avatar, m.rating, m.is_free 
							from movies as m join user_watchlist as uw on m.id = uw.movie_id
							where uw.user_id = $1;`

func GetWatchlist(userId uint, db *database.DBManager) ([]domain.Movie, error) {
	data, err := db.Query(querySelectUserWatchlist, userId)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Movie, 0)
	for _, movie := range data {
		result = append(result, domain.Movie{
			ID:     cast.ToUint(movie[0]),
			Title:  cast.ToString(movie[1]),
			Avatar: cast.ToString(movie[2]),
			Rating: cast.ToFloat(movie[3]),
			IsFree: cast.ToBool(movie[4]),
		})
	}
	return result, nil
}

func SendWatchList(movies []domain.Movie, sendTo string) error {
	message := FormFilmList(movies)
	SendEmailGomail(message, sendTo)
	return nil
}

func SendWatchlists(db *database.DBManager) error {
	users, err := GetSendInfo(db)
	if err != nil {
		return err
	}

	users = []domain.User{
		domain.User{
			ID:    1,
			Email: "pavel.cheklin@yandex.ru",
		},
		//	domain.User{
		//		ID:    2,
		//		Email: "tscheklin@gmail.com",
		//	},
	}

	for _, user := range users {
		watchlist, err := GetWatchlist(user.ID, db)
		if err != nil {
			log.Println(err)
		}
		err = SendWatchList(watchlist, user.Email)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func main() {
	db := database.Connect()
	err := SendWatchlists(db)
	if err != nil {
		log.Fatal(err)
	}
}
