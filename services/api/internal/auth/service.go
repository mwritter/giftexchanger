package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidToken = errors.New("invalid or expired magic link")
	ErrUnauthorized = errors.New("unauthorized")
)

type Config struct {
	BaseURL      string
	MagicLinkTTL time.Duration
	SessionTTL   time.Duration
	CookieSecure bool
	CookieName   string
}

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = "http://localhost:3000"
	}
	if c.MagicLinkTTL == 0 {
		c.MagicLinkTTL = 15 * time.Minute
	}
	if c.SessionTTL == 0 {
		c.SessionTTL = 30 * 24 * time.Hour
	}
	if c.CookieName == "" {
		c.CookieName = "giftswap_session"
	}
	return c
}

type User struct {
	ID          string
	Email       string
	DisplayName string
}

type Service struct {
	pool   *pgxpool.Pool
	mailer Mailer
	cfg    Config
}

func NewService(pool *pgxpool.Pool, mailer Mailer, cfg Config) *Service {
	return &Service{pool: pool, mailer: mailer, cfg: cfg.withDefaults()}
}

func (s *Service) CookieName() string        { return s.cfg.CookieName }
func (s *Service) CookieSecure() bool        { return s.cfg.CookieSecure }
func (s *Service) SessionTTL() time.Duration { return s.cfg.SessionTTL }

func NormalizeEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" {
		return "", ErrInvalidEmail
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", ErrInvalidEmail
	}
	return email, nil
}

func (s *Service) RequestMagicLink(ctx context.Context, rawEmail string) error {
	email, err := NormalizeEmail(rawEmail)
	if err != nil {
		return err
	}

	raw, err := newRandomToken()
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}
	hash := HashToken(raw)
	expires := time.Now().Add(s.cfg.MagicLinkTTL)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE magic_link_tokens
		SET consumed_at = now()
		WHERE lower(email) = $1 AND consumed_at IS NULL
	`, email); err != nil {
		return fmt.Errorf("invalidate prior tokens: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO magic_link_tokens (email, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, email, hash, expires); err != nil {
		return fmt.Errorf("insert token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	link := fmt.Sprintf("%s/api/auth/callback?token=%s", strings.TrimRight(s.cfg.BaseURL, "/"), raw)
	if err := s.mailer.SendMagicLink(ctx, email, link); err != nil {
		return fmt.Errorf("send magic link: %w", err)
	}
	return nil
}

func (s *Service) ConsumeMagicLink(ctx context.Context, rawToken string) (sessionRaw string, expiresAt time.Time, err error) {
	if strings.TrimSpace(rawToken) == "" {
		return "", time.Time{}, ErrInvalidToken
	}

	hash := HashToken(rawToken)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx)

	var email string
	err = tx.QueryRow(ctx, `
		UPDATE magic_link_tokens
		SET consumed_at = now()
		WHERE token_hash = $1
		  AND consumed_at IS NULL
		  AND expires_at > now()
		RETURNING email
	`, hash).Scan(&email)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", time.Time{}, ErrInvalidToken
	}
	if err != nil {
		return "", time.Time{}, fmt.Errorf("consume token: %w", err)
	}

	var userID string
	err = tx.QueryRow(ctx, `
		INSERT INTO users (email)
		VALUES ($1)
		ON CONFLICT ((lower(email))) DO UPDATE SET email = users.email
		RETURNING id
	`, email).Scan(&userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("upsert user: %w", err)
	}

	sessionRaw, err = newRandomToken()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("generate session: %w", err)
	}
	expiresAt = time.Now().Add(s.cfg.SessionTTL)
	if _, err := tx.Exec(ctx, `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, HashToken(sessionRaw), expiresAt); err != nil {
		return "", time.Time{}, fmt.Errorf("insert session: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", time.Time{}, fmt.Errorf("commit: %w", err)
	}
	return sessionRaw, expiresAt, nil
}

func (s *Service) Logout(ctx context.Context, sessionRaw string) error {
	if sessionRaw == "" {
		return nil
	}
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, HashToken(sessionRaw))
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (s *Service) UserFromSession(ctx context.Context, sessionRaw string) (*User, error) {
	if sessionRaw == "" {
		return nil, ErrUnauthorized
	}

	var user User
	err := s.pool.QueryRow(ctx, `
		SELECT u.id::text, u.email, u.display_name
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > now()
	`, HashToken(sessionRaw)).Scan(&user.ID, &user.Email, &user.DisplayName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, fmt.Errorf("lookup session: %w", err)
	}
	return &user, nil
}
