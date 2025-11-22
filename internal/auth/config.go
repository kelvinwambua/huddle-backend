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
	SessionName  = "huddle_session"
	SessionIDKey = "session_id"
	MaxAge       = 86400 * 7 // 7 days
)

var Store *sessions.CookieStore

func InitAuth() {
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		panic("SESSION_SECRET must be set in production")
	}

	appEnv := os.Getenv("APP_ENV")
	isProduction := appEnv == "production"

	Store = sessions.NewCookieStore([]byte(sessionSecret))
	Store.MaxAge(MaxAge)
	Store.Options.Path = "/"
	Store.Options.HttpOnly = true
	Store.Options.Secure = isProduction
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
