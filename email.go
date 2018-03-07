package bart

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os/exec"

	"github.com/hoisie/mustache"
	"github.com/howeyc/gopass"
)

type authPair struct {
	login    string
	password string
}

// Asks user for the login and password for the email server
func (ap *authPair) prompt() error {
	fmt.Printf("Login: ")
	fmt.Scanln(&ap.login)

	fmt.Printf("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return err
	}
	ap.password = string(pass) // `pass` is already a slice

	return nil
}

type EmailBuilder interface {
	AddAuthor(*Author) EmailBuilder
	AddRecipient(string) EmailBuilder
	AddContent(string) EmailBuilder
	Build(map[string]string) Email
}

// Email builder struct containing all the nitty-gritty details and subfields of our custom email
type emailBuilder struct {
	fromName  string
	fromEmail string
	toEmail   string
	mailText  string
}

func NewEmail() EmailBuilder {
	return &emailBuilder{}
}

func (eb *emailBuilder) AddAuthor(a *Author) EmailBuilder {
	eb.fromName = a.Name
	eb.fromEmail = a.Email

	return eb
}

func (eb *emailBuilder) AddRecipient(rs string) EmailBuilder {
	eb.toEmail = rs

	return eb
}

func (eb *emailBuilder) AddContent(s string) EmailBuilder {
	eb.mailText = s

	return eb
}

func (eb *emailBuilder) Build(context map[string]string) Email {
	body := "From: " + EncodeRfc1342(eb.fromName) + "<" + eb.fromEmail + ">\r\n"
	body += "To: " + eb.toEmail + "\r\n"
	body += "Subject: {{__subject_encoded__}}\r\n"
	body += "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body += eb.mailText
	body = mustache.Render(body, context)

	e := new(email)
	e.fromEmail = eb.fromEmail
	e.toEmails = []string{eb.toEmail, eb.fromEmail}
	e.body = []byte(body)

	return e
}

type Email interface {
	Send(*EmailServer, *authPair) error
	OpenInBrowser(string) error
	GetRecipients() []string
}

type email struct {
	fromEmail string
	toEmails  []string
	body      []byte
}

func (e *email) Send(s *EmailServer, ap *authPair) error {
	serverInfo := fmt.Sprintf("%s:%d", s.Hostname, s.Port)
	auth := smtp.PlainAuth("", ap.login, ap.password, s.Hostname)

	return smtp.SendMail(serverInfo, auth, e.fromEmail, e.toEmails, e.body)
}

func (e *email) OpenInBrowser(browserName string) error {
	html := append(
		[]byte(`<html><head><meta charset="UTF-8"></head>`),
		e.body...,
	)
	html = append(html, []byte(`</html>`)...)

	tmpfile, err := ioutil.TempFile(".", "bart_preview_")
	if err != nil {
		return err
	}
	if _, err := tmpfile.Write(html); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	cmd := exec.Command(browserName, tmpfile.Name())
	return cmd.Start()
}

func (e *email) GetRecipients() []string {
	return e.toEmails
}
