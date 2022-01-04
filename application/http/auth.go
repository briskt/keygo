package http

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/google"
)

const AuthCallbackPath = "/auth/callback"

type AuthError struct {
	Error string
}

type ProviderOption struct {
	Key         string
	Name        string
	RedirectURL string
}

type providerConfig struct {
	key         string
	name        string
	authKey     string
	secret      string
	callbackURL string
	domain      string
}

var providers = []providerConfig{
	{
		key:         "auth0",
		name:        "Auth0",
		authKey:     os.Getenv("AUTH0_KEY"),
		secret:      os.Getenv("AUTH0_SECRET"),
		callbackURL: os.Getenv("HOST") + AuthCallbackPath + "?provider=auth0",
		domain:      os.Getenv("AUTH0_DOMAIN"),
	}, {
		key:         "google",
		name:        "Google",
		authKey:     os.Getenv("GOOGLE_KEY"),
		secret:      os.Getenv("GOOGLE_SECRET"),
		callbackURL: os.Getenv("HOST") + AuthCallbackPath + "?provider=google",
	},
}

func init() {
	for _, p := range providers {
		if p.secret == "" || p.authKey == "" {
			continue
		}
		switch p.key {
		case "auth0":
			goth.UseProviders(auth0.New(p.authKey, p.secret, p.callbackURL, p.domain))
		case "google":
			goth.UseProviders(google.New(p.authKey, p.secret, p.callbackURL))
		}
	}
}

func RegisterAuthRoutes(e *echo.Echo) {
	e.GET("/auth/login", authLogin)
	e.GET(AuthCallbackPath, authCallback)
	e.GET("/auth/logout", authLogout)
}

func authLogin(c echo.Context) error {
	options := make([]ProviderOption, 0)
	for _, p := range providers {
		if p.secret == "" || p.authKey == "" {
			continue
		}
		gothic.GetProviderName = func(req *http.Request) (string, error) { return p.key, nil }
		url, err := gothic.GetAuthURL(c.Response(), c.Request())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
		}
		options = append(options, ProviderOption{
			Key:         p.key,
			Name:        p.name,
			RedirectURL: url,
		})
	}
	return c.JSON(http.StatusOK, options)
}

func authLogout(c echo.Context) error {
	err := gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusTemporaryRedirect, os.Getenv("UI_URL"))
}

func authCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}
