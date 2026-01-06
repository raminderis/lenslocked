package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main2() {
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
	from := "test@lenslocked.com"
	to := "raminderis@live.com"
	subject := "This is a test email subject line"
	plaintext := "This is a test email body"
	html := `<h1> Hello There buddy! </h1><p>This is the email</p><Hope you are enjoy it</p>`
	msg := gomail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("test/plain", plaintext)
	msg.AddAlternative("text/html", html)
	msg.WriteTo(os.Stdout)
	fmt.Println("moving to send the mail")
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.SSL = false
	err = dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("mail sent successfully")
}
