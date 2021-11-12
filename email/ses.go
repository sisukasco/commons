package email

import (
	"bytes"
	"log"

	"github.com/sisukas/commons/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

// credits : https://github.com/tj/go-ses

func awsStrings(arr []string) []*string {
	ret := []*string{}
	for _, a := range arr {
		ret = append(ret, aws.String(a))
	}
	return ret
}

func createSession(ses *SESCredentials) (*session.Session, error) {
	cred := credentials.NewStaticCredentials(
		ses.AccessKey,
		ses.SecretKey, "")
	if cred == nil {
		return nil, errors.New("Error loading AWS SES credentials for emailing")
	}

	awsConfig := &aws.Config{
		Region:      &ses.AWSRegion,
		Credentials: cred,
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating SES session")
	}

	return sess, nil
}

// *********************************************************************
//	create and send text or html email to single receipents.
//	@returns resp *ses.SendEmailOutput
//
func sendEmailSES(sesCred *SESCredentials, emailData *Email) error {

	// start a new aws session
	sess, err := createSession(sesCred)
	if err != nil {
		return err
	}
	// start a new ses session
	svc := ses.New(sess)

	msg := gomail.NewMessage()

	msg.SetHeader("From", emailData.From)
	msg.SetHeader("To", emailData.To...)
	msg.SetHeader("Cc", emailData.Cc...)
	msg.SetHeader("Bcc", emailData.Bcc...)
	msg.SetHeader("Subject", emailData.Subject)
	msg.SetBody("text/html", emailData.BodyHTML)

	for _, att := range emailData.Attachments {
		msg.Attach(att.Path, gomail.Rename(att.FileName))
	}

	if len(emailData.ReplyTo) > 1 {
		msg.SetHeader("Reply-To", emailData.ReplyTo)
	}
	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	bmsg := &ses.RawMessage{Data: emailRaw.Bytes()}

	dests := emailData.To
	dests = append(dests, emailData.Cc...)
	dests = append(dests, emailData.Bcc...)

	params := &ses.SendRawEmailInput{
		Destinations: awsStrings(dests),
		RawMessage:   bmsg,
		Source:       aws.String(emailData.From), // Required
	}

	// send email
	resp, err := svc.SendRawEmail(params)

	if err != nil {
		log.Printf("SES error %v ", err)
		return errors.Wrap(err, "sendEmailSES")
	}

	log.Printf(" SES Response %v ", utils.ToJSONString(resp))
	return nil
}
