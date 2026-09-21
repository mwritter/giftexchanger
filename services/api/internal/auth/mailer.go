package auth

import (
	"context"
	"log"
)

type Mailer interface {
	SendMagicLink(ctx context.Context, to, url string) error
}

// LogMailer writes the login URL to the process log. Used in local development
// so we do not need a real inbox.
type LogMailer struct {
	Logger *log.Logger
}

func (m LogMailer) SendMagicLink(_ context.Context, to, url string) error {
	logger := m.Logger
	if logger == nil {
		logger = log.Default()
	}
	logger.Printf("magic link for %s: %s", to, url)
	return nil
}
