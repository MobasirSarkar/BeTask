package auth

import "github.com/gorilla/sessions"

const (
	SessionName = "session"
)

type SessionsOptions struct {
	CookiesKey string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
}

func NewSessionStore(opts SessionsOptions) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(opts.CookiesKey))
	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure

	return store
}
