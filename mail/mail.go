package mail

import (
	"os"
	"fmt"
	"gopkg.in/gomail.v2"
)

func SendMail(name string, to string, authCode string) error {

	from := os.Getenv("EMAIL_ID")
	password := os.Getenv("EMAIL_PASS")

	// 	// smtp server configuration.
	smtpHost := "smtp-mail.outlook.com"
	smtpPort := 587

	subject := "Puppy Love Authentication Code"
	body := "<!DOCTYPE html><html><style>.container{font-size: large;}body{font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;}code{padding: 6px;background-color: antiquewhite;border-radius: 5px;}</style><body><div class='container'>Hello" + name + ",\nThis is your <code>"+ authCode+ "</code> for Puppy Love.</div></body></html>"

	mail_server := gomail.NewMessage()
	mail_server.SetHeader("From", from)
	mail_server.SetHeader("To", to)
	mail_server.SetHeader("Subject", subject)
	mail_server.SetBody("text/plain", body)

	client := gomail.NewDialer(smtpHost, smtpPort, from, password)
	if err := client.DialAndSend(mail_server); err != nil {
		fmt.Println("Error Sending Mail: ", err)
		return err
	}
	
	return nil
}