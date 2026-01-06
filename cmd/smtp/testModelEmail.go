package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/raminderis/lenslocked/models"
)

func main() {
	// email := models.Email{
	// 	From:      "test@lenslocked.com",
	// 	To:        "raminderis@live.com",
	// 	Subject:   "This is a test email subject line",
	// 	Plaintext: "This is a test email body",
	// 	HTML:      `<h1> Hello There buddy! </h1><p>This is the email</p><Hope you are enjoy it</p>`,
	// }
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtp_ssl_string := os.Getenv("SMTP_SSL")
	smtp_ssl_bool, err := strconv.ParseBool(smtp_ssl_string)
	if err != nil {
		panic(err)
	}
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	es.DailerSSL = smtp_ssl_bool
	// err := es.Send(email)
	// if err != nil {
	// 	panic(err)
	// }
	err = es.ForgotPassword("ramidner3@live.com", "https://lenslockex.com/reset-pw?token=123abc")
	if err != nil {
		panic(err)
	}
	fmt.Println("mail sent successfully")
}
