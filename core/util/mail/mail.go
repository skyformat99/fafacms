package mail

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
)

var Debug = false

type Sender struct {
	//Bcc      string
	//BccName  string
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Subject  string `json:"Subject"`
	Body     string `json:"Body"`
}

type Message struct {
	To     string
	ToName string
	Sender
}

func (mm *Message) Sent() error {
	if Debug {
		return nil
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mm.Email, "FaFaCMS")
	m.SetAddressHeader("To", mm.To, mm.ToName)
	m.SetHeader("Subject", mm.Subject)

	//m.SetHeader("Bcc",
	//	m.FormatAddress(mm.Bcc, mm.BccName))

	m.SetBody("text/html", mm.Body)

	d := gomail.NewDialer(mm.Host, mm.Port, mm.Email, mm.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
