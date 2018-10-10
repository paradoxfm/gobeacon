package service

import (
	"bytes"
	"html/template"
	"net/smtp"
)

const (
	Subject = "Cброс пароля"

	SMTP_SERVER = "smtp.gmail.com"
	User        = "ControlCenter66@gmail.com"
	Password    = "MJe#NB673B5&"
)

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

var auth smtp.Auth

func sendNewPassword(email, password string) (bool, error) {
	auth = smtp.PlainAuth("", User, Password, SMTP_SERVER)
	templateData := struct{ Password string }{Password: password,}
	r := NewRequest([]string{email}, Subject, "Сброс пароля!")
	err := r.parseTemplate("assets/restorePassword.html", templateData)
	if err == nil {
		return r.sendEmail()
	}
	return false, err
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) sendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, User, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
