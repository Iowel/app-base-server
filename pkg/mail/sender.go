package mail

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	Sendmail(
		tmpl string,
		data interface{},
		subject string,
		to string,
		cc, bcc, attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

//go:embed templates
var emailTemplateFS embed.FS

func (sender *GmailSender) Sendmail(
	tmpl string,
	data interface{},
	subject string,
	to string,
	cc, bcc, attachFiles []string,
) error {
	// HTML шаблон
	htmlTemplatePath := fmt.Sprintf("templates/%s.html.tmpl", tmpl)
	tHtml, err := template.New("email-html").ParseFS(emailTemplateFS, htmlTemplatePath)
	if err != nil {
		log.Println("Ошибка при парсинге HTML шаблона:", err)
		return err
	}

	var htmlBuf bytes.Buffer
	if err := tHtml.ExecuteTemplate(&htmlBuf, "body", data); err != nil {
		log.Println("Ошибка при рендеринге HTML:", err)
		return err
	}

	// Plain шаблон
	plainTemplatePath := fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	tPlain, err := template.New("email-plain").ParseFS(emailTemplateFS, plainTemplatePath)
	if err != nil {
		log.Println("Ошибка при парсинге plain шаблона:", err)
		return err
	}

	var plainBuf bytes.Buffer
	if err := tPlain.ExecuteTemplate(&plainBuf, "body", data); err != nil {
		log.Println("Ошибка при рендеринге plain текста:", err)
		return err
	}


	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.To = []string{to}
	e.Cc = cc
	e.Bcc = bcc
	e.Subject = subject
	e.HTML = htmlBuf.Bytes()
	e.Text = plainBuf.Bytes()

	for _, f := range attachFiles {
		if _, err := e.AttachFile(f); err != nil {
			return fmt.Errorf("не удалось прикрепить файл %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
