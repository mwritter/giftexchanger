package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/mwritter/giftswap/services/api/internal/auth"
)

type magicLinkRequest struct {
	Email string `json:"email"`
}

func (s *Server) requestMagicLink(w http.ResponseWriter, r *http.Request) {
	var req magicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err := s.auth.RequestMagicLink(r.Context(), req.Email)
	if errors.Is(err, auth.ErrInvalidEmail) {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if err != nil {
		log.Printf("request magic link: %v", err)
		writeError(w, http.StatusInternalServerError, "could not send magic link")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) authCallback(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	sessionRaw, expiresAt, err := s.auth.ConsumeMagicLink(r.Context(), token)
	if errors.Is(err, auth.ErrInvalidToken) {
		writeError(w, http.StatusBadRequest, "invalid or expired magic link")
		return
	}
	if err != nil {
		log.Printf("consume magic link: %v", err)
		writeError(w, http.StatusInternalServerError, "could not complete login")
		return
	}

	s.setSessionCookie(w, sessionRaw, expiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(s.auth.CookieName())
	raw := ""
	if cookie != nil {
		raw = cookie.Value
	}
	if err := s.auth.Logout(r.Context(), raw); err != nil {
		log.Printf("logout: %v", err)
		writeError(w, http.StatusInternalServerError, "could not log out")
		return
	}
	s.clearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"id":           user.ID,
		"email":        user.Email,
		"display_name": user.DisplayName,
	})
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(s.auth.CookieName())
		if err != nil || cookie.Value == "" {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		user, err := s.auth.UserFromSession(r.Context(), cookie.Value)
		if errors.Is(err, auth.ErrUnauthorized) {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		if err != nil {
			log.Printf("load session: %v", err)
			writeError(w, http.StatusInternalServerError, "could not load session")
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
	})
}

func (s *Server) setSessionCookie(w http.ResponseWriter, raw string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.auth.CookieName(),
		Value:    raw,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   s.auth.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Server) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.auth.CookieName(),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.auth.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
	})
}
