package email

import (
	"log"
	"mime"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/cbroglie/mustache"
	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
)

type SMTPLogin struct {
	Host     string
	UserName string
	Password string
}

type SESCredentials struct {
	AWSRegion string
	AccessKey string
	SecretKey string
}
type EmailConfig struct {
	From           string
	SMTPLogin      *SMTPLogin
	SEScredentials *SESCredentials
}

func NewSMTPConfig(from string, host string, user string, pass string) *EmailConfig {
	smtp := &SMTPLogin{Host: host, UserName: user, Password: pass}
	return &EmailConfig{From: from, SMTPLogin: smtp}
}

func NewSESConfig(from string, region string, accessKey string, secretKey string) *EmailConfig {
	ses := &SESCredentials{AWSRegion: region, AccessKey: accessKey, SecretKey: secretKey}
	return &EmailConfig{From: from, SEScredentials: ses}
}

type Emailer struct {
	conf *EmailConfig
}

func NewEmailer(conf *EmailConfig) *Emailer {

	return &Emailer{conf: conf}
}

func (emx *Emailer) SendEmail(to string, subject string, body string) error {

	htmlEmailBody, err := mustache.Render(htmlEmailTemplate, map[string]string{
		"Subject": subject,
		"Body":    body,
	})
	if err != nil {
		return errors.Wrap(err, "making email template")
	}

	if emx.conf.SMTPLogin != nil {
		emx.sendBySMTP(to, subject, htmlEmailBody)
	} else if emx.conf.SEScredentials != nil {
		emx.sendBySES(to, subject, htmlEmailBody)
	} else {
		log.Printf("sending email to %s subject %s", to, subject)
	}

	return nil
}

func (emx *Emailer) SendEmailAdv(em *Email) error {

	htmlEmailBody, err := mustache.Render(htmlEmailTemplate, map[string]string{
		"Subject": em.Subject,
		"Body":    em.BodyHTML,
	})
	if err != nil {
		return errors.Wrap(err, "making email template")
	}
	em.BodyHTML = htmlEmailBody

	if len(em.From) <= 0 {
		em.From = emx.conf.From
	}

	if emx.conf.SMTPLogin != nil {
		return emx.sendBySMTPAdv(em)
	} else if emx.conf.SEScredentials != nil {
		return sendEmailSES(emx.conf.SEScredentials, em)
	} else {
		log.Printf("sending email to %v subject %s", em.To, em.Subject)
	}

	return nil
}

//TODO: use https://github.com/jaytaylor/html2text to convert from HTML to text

func (em *Emailer) attachFile(e *email.Email, filename string, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	ct := mime.TypeByExtension(filepath.Ext(path))
	basename := filepath.Base(filename)
	_, err = e.Attach(f, basename, ct)
	return err
}

func (emx *Emailer) sendBySMTPAdv(em *Email) error {

	e := email.NewEmail()
	e.From = em.From
	if len(e.From) <= 0 {
		e.From = emx.conf.From
	}
	e.To = em.To
	e.Cc = em.Cc
	e.Bcc = em.Bcc
	e.Subject = em.Subject
	e.ReplyTo = []string{em.ReplyTo}
	e.HTML = []byte(em.BodyHTML)

	for _, att := range em.Attachments {
		err := emx.attachFile(e, att.FileName, att.Path)
		if err != nil {
			log.Printf("Error attaching file to email %v", err)
		}
	}
	smtpServer := emx.conf.SMTPLogin.Host
	user := emx.conf.SMTPLogin.UserName
	pass := emx.conf.SMTPLogin.Password

	//localhost because PlainAuth will only send the credentials if the connection is using TLS
	// or is connected to localhost
	err := e.Send(smtpServer, smtp.PlainAuth("", user, pass, "localhost"))

	if err != nil {
		log.Printf("SMTP error %v", err)
		return errors.Wrap(err, "sendBySMTP")
	}
	return nil
}

func (emx *Emailer) sendBySES(to string, subject string, body string) error {
	e := &Email{}

	if len(e.From) <= 0 {
		e.From = emx.conf.From
	}
	e.To = []string{to}
	e.Subject = subject
	e.BodyHTML = body

	return sendEmailSES(emx.conf.SEScredentials, e)
}

//TODO: Add ctx context
func (emx *Emailer) sendBySMTP(to string, subject string, body string) error {

	e := email.NewEmail()

	e.From = emx.conf.From
	e.To = []string{to} //[]string{"test@example.com"}
	//e.Bcc = []string{"test_bcc@example.com"}
	//e.Cc = []string{"test_cc@example.com"}
	e.Subject = subject
	//e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte(body)

	smtpServer := emx.conf.SMTPLogin.Host
	pass := emx.conf.SMTPLogin.Password

	err := e.Send(smtpServer, smtp.PlainAuth("", e.From, pass, "localhost"))

	if err != nil {
		log.Printf("SMTP error %v", err)
		return errors.Wrap(err, "sendBySMTP")
	}

	return nil
}
