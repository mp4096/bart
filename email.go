package bart

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cbroglie/mustache"
	"github.com/howeyc/gopass"
	"github.com/jaytaylor/html2text"
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
	Build(map[string]string) (Email, error)
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

func (eb *emailBuilder) Build(context map[string]string) (Email, error) {
	mustache.AllowMissingVariables = false

	headerTemplate := "From: " + EncodeRfc1342(eb.fromName) + " <" + eb.fromEmail + ">\r\n" +
		"To: " + eb.toEmail + "\r\n" +
		"Subject: {{__subject_encoded__}}\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: multipart/alternative;\r\n" +
		"\tboundary=\"" + boundary + "\"\r\n\r\n"

	header, err := mustache.Render(headerTemplate, context)
	if err != nil {
		return nil, err
	}

	partHtml, err := mustache.Render(eb.mailText, context)
	if err != nil {
		return nil, err
	}

	partText, err := html2text.FromString(partHtml, html2text.Options{PrettyTables: true})
	if err != nil {
		return nil, err
	}

	e := new(email)
	e.fromEmail = eb.fromEmail
	e.toEmails = []string{eb.toEmail, eb.fromEmail}
	e.header = emailHeader{header}
	e.items = []mimepart{&textHtml{partHtml}, &textPlain{partText}}

	return e, nil
}

type Email interface {
	Send(*EmailServer, *authPair) error
	OpenInBrowser(string) error
	GetRecipients() []string
}

type email struct {
	fromEmail string
	toEmails  []string
	header    emailHeader
	items     []mimepart
}

func (e *email) Send(s *EmailServer, ap *authPair) error {
	serverInfo := fmt.Sprintf("%s:%d", s.Hostname, s.Port)
	auth := smtp.PlainAuth("", ap.login, ap.password, s.Hostname)

	m := e.header.asPlainBytes()
	for i := 0; i < len(e.items); i++ {
		m = append(m, e.items[i].asBase64()...)
	}
	m = append(m, footerAsPlainBytes()...)
	return smtp.SendMail(serverInfo, auth, e.fromEmail, e.toEmails, m)
}

func (e *email) OpenInBrowser(browserName string) error {
	// Create an HTML view of the email
	html := e.header.asHtml()
	for i := 0; i < len(e.items); i++ {
		html = append(html, e.items[i].asHtml()...)
	}
	html = append(html, footerAsHtml()...)

	// Create a temp file and write the email to it
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

	// Rename the file
	oldFilename, err := filepath.Abs(filepath.Join(".", tmpfile.Name()))
	if err != nil {
		return err
	}
	newFilename := oldFilename + ".html"
	os.Rename(oldFilename, newFilename)

	// Open it in browser
	cmd := exec.Command(browserName, newFilename)
	return cmd.Start()
}

func (e *email) GetRecipients() []string {
	return e.toEmails
}
