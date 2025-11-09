package providing

import (
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/emailing"
)

func NewEmailing(cfg config.Config) emailing.Emailing {
	return emailing.NewService(
		cfg.Emailing.SmtpHost,
		cfg.Emailing.SmtpPort,
		cfg.Emailing.SmtpUser,
		cfg.Emailing.SmtpPassword,
		cfg.Emailing.DefaultFromEmail,
		cfg.Emailing.DefaultFromTitle,
	)
}
