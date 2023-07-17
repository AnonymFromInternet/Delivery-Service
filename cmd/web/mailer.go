package main

import (
	"bytes"
	"fmt"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"sync"
	"time"
)

type MailServer struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Wait        *sync.WaitGroup
	MailerChan  chan Message
	ErrorChan   chan error
	DoneChan    chan bool
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
	Template    string
}

// a function to listen for messages on the MailerChan

func (m *MailServer) sendMail(msg Message, errorChan chan error) {
	if msg.Template == "" {
		msg.Template = "mail"
	}

	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	// build html version
	htmlMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		errorChan <- err
	}

	// build plain text version
	plainTextMessage, err := m.buildPlaynTextMessage(msg)
	if err != nil {
		errorChan <- err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 5 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainTextMessage)
	email.AddAlternative(mail.TextHTML, htmlMessage)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		errorChan <- err
	}
}

func (m *MailServer) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.html.gohtml", msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", msg.DataMap)
	if err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	return m.inlineCSS(formattedMessage)
}

func (m *MailServer) buildPlaynTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.plain.gohtml", msg.Template)

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", msg.DataMap)
	if err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *MailServer) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	_premailer, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := _premailer.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *MailServer) getEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
