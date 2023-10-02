package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(name string, to string, authCode string) error {

	from := os.Getenv("EMAILID")
	password := os.Getenv("EMAILPASS")

	// 	// smtp server configuration.
	smtpHost := "smtp.cc.iitk.ac.in"
	smtpPort := "25"

	var toSend []string
	toSend = append(toSend, to)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("From:PClub IITK\nSubject: Puppy Love Authentication Code \n%s\n\n", mimeHeaders)))

	mailTemplate := fmt.Sprintf("<!DOCTYPE html><html><style>.container{font-size: large;}body{font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;}code{padding: 6px;background-color: antiquewhite;border-radius: 5px;}</style><body><div class='container'>Hello %s,\nThis is your <code>%s</code> for Puppy Love.</div></body></html>", name, authCode)
	body.Write([]byte(mailTemplate))

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toSend, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}