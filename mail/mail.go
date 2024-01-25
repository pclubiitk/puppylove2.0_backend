package mail

import (
	"os"
	"fmt"
	"gopkg.in/gomail.v2"
)

func SendMail(name string, to string, authCode string) error {

	/* Old Mail function for IIT K Mail Servers
	from := os.Getenv("EMAIL_ID")
	password := os.Getenv("EMAIL_PASS")

	// 	// smtp server configuration.
	smtpHost := "mmtp.iitk.ac.in"
	smtpPort := "25"

	var toSend []string
	toSend = append(toSend, to)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("From:PClub IITK <pclubiitk@gmail.com>\nSubject: Puppy Love Authentication Code \n%s\n\n", mimeHeaders)))

	mailTemplate := fmt.Sprintf("<!DOCTYPE html><html><style>.container{font-size: large;}body{font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;}code{padding: 6px;background-color: antiquewhite;border-radius: 5px;}</style><body><div class='container'>Hello %s,\nThis is your <code>%s</code> for Puppy Love.</div></body></html>", name, authCode)
	body.Write([]byte(mailTemplate))

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toSend, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
	*/

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