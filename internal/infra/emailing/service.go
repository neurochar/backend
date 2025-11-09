package emailing

import (
	"io"

	"github.com/k3a/html2text"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"gopkg.in/gomail.v2"
)

type service struct {
	smtpHost         string
	smtpPort         int
	smtpUser         string
	smtpPassword     string
	defaultFromEmail string
	defaultFromTitle string
}

func NewService(
	smtpHost string,
	smtpPort int,
	smtpUser string,
	smtpPassword string,
	defaultFromEmail string,
	defaultFromTitle string,
) *service {
	return &service{
		smtpHost:         smtpHost,
		smtpPort:         smtpPort,
		smtpUser:         smtpUser,
		smtpPassword:     smtpPassword,
		defaultFromEmail: defaultFromEmail,
		defaultFromTitle: defaultFromTitle,
	}
}

func (s *service) prepareMsg(msg Message) (*gomail.Message, error) {
	if len(msg.To) == 0 {
		return nil, ErrIncorrectTo
	}

	if len(msg.Subject) == 0 {
		return nil, ErrIncorrectSubject
	}

	gomailMsg := gomail.NewMessage()
	gomailMsg.SetHeader("From", gomailMsg.FormatAddress(s.defaultFromEmail, s.defaultFromTitle))
	gomailMsg.SetHeader("To", msg.To)
	gomailMsg.SetHeader("Subject", msg.Subject)

	switch {
	case len(msg.TextPlain) > 0:
		gomailMsg.AddAlternative("text/plain", msg.TextPlain)

	case len(msg.TextHtml) > 0 && msg.AutoTextPlain:
		gomailMsg.AddAlternative("text/plain", html2text.HTML2Text(msg.TextHtml))
	}

	if len(msg.TextHtml) > 0 {
		gomailMsg.AddAlternative("text/html", msg.TextHtml)
	}

	for _, file := range msg.Files {
		gomailMsg.Attach(file.Name,
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(file.Data)
				return err
			}),
			gomail.SetHeader(file.Headers),
		)
	}

	return gomailMsg, nil
}

func (s *service) sendMsg(gomailMsg *gomail.Message) error {
	dialer := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPassword)

	if err := dialer.DialAndSend(gomailMsg); err != nil {
		return appErrors.ErrInternal.WithWrap(err)
	}

	return nil
}

func (s *service) Send(msg Message) error {
	gomailMsg, err := s.prepareMsg(msg)
	if err != nil {
		return err
	}

	err = s.sendMsg(gomailMsg)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendAsyc(msg Message, onSendError func(error)) error {
	gomailMsg, err := s.prepareMsg(msg)
	if err != nil {
		return err
	}

	go func() {
		err = s.sendMsg(gomailMsg)
		if err != nil && onSendError != nil {
			onSendError(err)
		}
	}()

	return nil
}
