package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

const (
	SessionName = "huddle_auth_session"
	MaxAge      = 86400 * 30
)

var Store *sessions.CookieStore

func InitAuth() {
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "default-secret-please-change-in-production"
	}

	Store = sessions.NewCookieStore([]byte(sessionSecret))
	Store.MaxAge(MaxAge)
	Store.Options.Path = "/"
	Store.Options.HttpOnly = true
	Store.Options.Secure = false
	Store.Options.SameSite = 2

	gothic.Store = Store

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALLBACK_URL"),
		),
		github.New(
			os.Getenv("GITHUB_CLIENT_ID"),
			os.Getenv("GITHUB_CLIENT_SECRET"),
			os.Getenv("GITHUB_CALLBACK_URL"),
		),
	)
}
