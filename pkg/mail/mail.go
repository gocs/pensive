package mail

import (
	"fmt"
	"net/smtp"

	"github.com/asaskevich/govalidator"
	"github.com/gocs/errored"
)

const (
	ErrNotValidFromAddr = errored.New("sender email address provided is invalid")
	ErrNotValidToAddr   = errored.New("recipient email address provided is invalid")
)

// Send this sends a mail specifically for gmail using plain auth
// fromAddr is the sender email address
// password is the sender password; you must use an App Password
// toAddr is the recipient email address
// subject is the subject of the mail
// body is the content of your mail
func Send(fromAddr, password, toAddr, subject, body string) error {
	if !govalidator.IsEmail(fromAddr) {
		return ErrNotValidFromAddr
	}

	if !govalidator.IsEmail(toAddr) {
		return ErrNotValidToAddr
	}

	msg := fmt.Sprint("From: ", fromAddr, "\n",
		"To: ", toAddr, "\n",
		"Subject: ", subject, "\n\n", body)

	return smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", fromAddr, password, "smtp.gmail.com"),
		fromAddr, []string{toAddr}, []byte(msg))
}
