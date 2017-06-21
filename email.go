package media

import (
	"os"
	"net/url"
	"crypto/md5"
	"fmt"
	"bytes"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
	"time"
	"html/template"
	"crypto/tls"
	"errors"
)

const TOKEN_LIFESPAN = 1 * time.Hour
const MAIL_FROM_ADDRESS = "noone@unittest.com"
const MAIL_SUBJECT = "Media login token"
const MAIL_REPLY_TO = "noone@unittest.com"
const MAIL_SERVER = "localhost"
const MAIL_USE_TLS = true

var mailTlsConfig = &tls.Config{

}

func sendTokenEmail(requestUrl *url.URL, emailAddress string) (err error) {
	// create the token
	token := &Token{
		Code: fmt.Sprintf("LoginToken-%x-%x-%x", GetRandomBytes(8), GetRandomBytes(12), GetRandomBytes(12)),
		Expiration: time.Now().Add(TOKEN_LIFESPAN),
		Email: emailAddress,
	}

	// token URL
	requestString := requestUrl.String()
	tokenUrl := requestString[0:strings.LastIndex(requestString, "/")] + "/" + PATH_TOKEN + "/" + token.Code

	// create a boundary ID
	boundary := fmt.Sprintf("%x", md5.Sum([]byte(token.Code)))

	// create the mail message
	emailFields := map[string]string{"LoginToken": tokenUrl}
	emailBody := bytes.NewBuffer(make([]byte, 0))

	// write RFC 822 message headers
	emailBody.WriteString(fmt.Sprintf("From: %s\n", MAIL_FROM_ADDRESS))
	emailBody.WriteString(fmt.Sprintf("To: %s\n", emailAddress))
	emailBody.WriteString(fmt.Sprintf("Subject: %s\n", MAIL_SUBJECT))
	emailBody.WriteString(fmt.Sprintf("Reply-to: %s\n", MAIL_REPLY_TO))

	// write extended headers
	emailBody.WriteString("MIME-Version: 1.0\n")
	emailBody.WriteString(fmt.Sprintf("Content-type: multipart/mixed; boundary=\"%s\"\n", boundary))
	emailBody.WriteString("X-Priority: 1\n")
	emailBody.WriteString("priority: Urgent\n")
	emailBody.WriteString("Importance: high\n")

	emailBody.WriteString("\n")        // blank line between subject and body
	emailBody.WriteString("> This message is in MIME format. Since your mail reader does not understand this format, some or all of this message may not be legible.\n\n")

	// write the HTML body
	emailBody.WriteString("--" + boundary + "\n")
	emailBody.WriteString("Content-type: text/html; charset=\"UTF-8\"\n")
	emailBody.WriteString("Content-transfer-encoding: quoted-printable\n\n\n")

	htmlWriter := quotedprintable.NewWriter(emailBody)
	htmlText, err := template.New("token_email").ParseFiles("templates/token_email.html")
	if err != nil {
		os.Stderr.WriteString("ERROR: couldn't parse token email file: " + err.Error() + "\n")
	}
	err = htmlText.Execute(htmlWriter, emailFields)
	if err != nil {
		os.Stderr.WriteString("ERROR: couldn't execute token email template: " + err.Error() + "\n")
	}
	err = htmlWriter.Close()
	if err != nil {
		os.Stderr.WriteString("ERROR: couldn't close HTML quoted-printable writer: " + err.Error() + "\n")
	}

	// queue the email
	conn, err := smtp.Dial(MAIL_SERVER)
	if err != nil {
		os.Stderr.WriteString("Couldn't connect to email server: " + err.Error() + "\n")
	}
	if MAIL_USE_TLS {
		err := conn.StartTLS(mailTlsConfig)
		if err != nil {
			os.Stderr.WriteString("Couldn't STARTTLS with email server: " + err.Error() + "\n")
		}
	}
	err = conn.Hello("EHLO")
	if err != nil {
		os.Stderr.WriteString("Problems introducing zuul-reset to SMTP server: " + err.Error() + "\n")
	}
	err = conn.Mail(MAIL_FROM_ADDRESS)
	if err != nil {
		os.Stderr.WriteString("Problems sending mail-from to SMTP server:" + err.Error() + "\n")
	}
	err = conn.Rcpt(token.Email)
	if err != nil {
		os.Stderr.WriteString("Problems sending rcpt-to to SMTP server: " + err.Error() + "\n")
	}
	w, err := conn.Data()
	if err != nil {
		os.Stderr.WriteString("Couldn't get data writer for mail message: " + err.Error() + "\n")
	}
	n, err := w.Write(emailBody.Bytes())
	if err != nil {
		os.Stderr.WriteString("Error writing mail message body: " + err.Error() + "\n")
	}
	if n != emailBody.Len() {
		os.Stderr.WriteString("Didn't flush entire message to mail server\n")
	}
	err = conn.Quit()
	if err != nil {
		// scan the error message -- 200s are OK!
		if !strings.HasPrefix(err.Error(), "2") {
			os.Stderr.WriteString("Error writing email message to SMTP server, may not have been sent: " + err.Error() + "\n")
		} else {
			os.Stderr.WriteString("Sent token reset message to vendor email " + token.Email + "\n")
			db := getTokenDatabase()
			db.StoreToken(token)
			err = nil
			return
		}
	}
	return
}

func processEmailToken(code string) (token *Token, err error) {
	// load the token from the database
	db := getTokenDatabase()
	token, err = db.GetToken(code)
	if err != nil {
		return
	}

	// is the token expired?
	if token.Expiration.Before(time.Now()) {
		return nil, errors.New("Token expired.")

	}

	// delete the token from the database
	err = db.DeleteToken(token)

	return
}

func getTokenDatabase() *Database {
	tokenDb, err := OpenDatabase(DB_TOKEN)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}

	return tokenDb
}

