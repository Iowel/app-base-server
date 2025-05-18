package mailer

import (
	"github.com/Iowel/app-base-server/configs"
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	conf *configs.Config
}

func NewMailer(conf *configs.Config) *Mailer {
	return &Mailer{
		conf: conf,
	}
}


//go:embed templates
var emailTemplateFS embed.FS

// отправляет email по шаблону, с данными и заданной темой
func (m *Mailer) SendMail(from, to, subject, tmpl string, data interface{}) error {
	// Формируем путь к HTML-шаблону
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	// парсим HTML-шаблон из встроенной файловой системы
	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		log.Println(err) 
		return err 
	}

	// буфер для хранения отрендеренного HTML-шаблона
	var tpl bytes.Buffer

	// Выполняем шаблон с переданными данными (например: имя, ссылка и т.д.)
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		log.Println(err)
		return err
	}

	// получаем HTML-контент как строку
	formattedMessage := tpl.String()

	// аналогично обрабатываем plain text шаблон (для текстовой версии письма)
	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		log.Println(err)
		return err
	}

	// выполняем plain шаблон и добавляем результат в тот же буфер
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		log.Println(err)
		return err
	}

	// получаем plain text содержимое
	plainMessage := tpl.String()

	// логируем обе версии письма — HTML и текстовую
	log.Println(formattedMessage, plainMessage)

	
	// Отправляем письмо

	// Настройка SMTP-сервера через библиотеку go-simple-mail
	server := mail.NewSMTPClient()
	server.Host = m.conf.Smtp.Host
	server.Port, _ = strconv.Atoi(m.conf.Smtp.Port)
	server.Username = m.conf.Smtp.Username
	server.Password = m.conf.Smtp.Password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Устанавливаем соединение с SMTP-сервером
	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// Создаем новое письмо
	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	// Устанавливаем HTML-контент как основной
	email.SetBody(mail.TextHTML, formattedMessage)

	// Добавляем альтернативную plain text версию (для клиентов без поддержки HTML)
	email.AddAlternative(mail.TextPlain, plainMessage)

	// Отправляем письмо через SMTP-клиент
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	// Логируем успешную отправку
	log.Println("send mail")

	return nil
}
