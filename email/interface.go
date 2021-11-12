package email

type Attachment struct {
	Path     string
	FileName string
}
type Email struct {
	ReplyTo     string
	From        string
	To          []string
	Bcc         []string
	Cc          []string
	Subject     string
	BodyHTML    string // Html message (optional)
	Attachments []Attachment
}

type EmailSender interface {
	SendEmail(to string, subject string, body string) error
	SendEmailAdv(e *Email) error
}

/*
type Email struct {
	From    string // From source email
	To      string // To destination email(s)
	Subject string // Subject text to send
	Text    string // Text is the text body representation
	HTML    string // HTMLBody is the HTML body representation
	ReplyTo string // Reply-To email(s)
}

*/
