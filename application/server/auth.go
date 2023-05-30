package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/google"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
)

const AuthCallbackPath = "/api/auth/callback"

const ClientIDParam = "client_id"

const ClientIDSessionKey = "ClientID"
const TokenSessionKey = "Token"

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

func (s *Server) authStatus(c echo.Context) error {
	var status keygo.AuthStatus

	token, err := s.getTokenFromSession(c)
	if err != nil {
		s.Logger.Errorf("error getting token from session: %s", err)
		return c.JSON(http.StatusOK, status)
	}

	if token.ExpiresAt.After(time.Now()) {
		status.IsAuthenticated = true
	}
	status.UserID = token.Auth.UserID
	status.Expiry = token.ExpiresAt

	return c.JSON(http.StatusOK, status)
}

func (s *Server) authLogin(c echo.Context) error {
	clientID := c.QueryParam(ClientIDParam)
	if clientID == "" {
		return c.JSON(http.StatusBadRequest, AuthError{Error: ClientIDParam + " is required to login"})
	}
	if err := sessionSetValue(c, ClientIDSessionKey, clientID); err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: "error saving clientID into session, " + err.Error()})
	}

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

func (s *Server) authLogout(c echo.Context) error {
	err := gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusTemporaryRedirect, os.Getenv("UI_URL"))
}

func (s *Server) authCallback(c echo.Context) error {
	clientID, err := sessionGetString(c, ClientIDSessionKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, AuthError{Error: ClientIDSessionKey + " not found in session"})
	}

	authUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())

	s.Logger.Infof("authCallback authUser = %+v", authUser)

	auth, err := s.AuthService.CreateAuth(c, keygo.Auth{
		Provider:   authUser.Provider,
		ProviderID: authUser.UserID,
		User: keygo.User{
			FirstName: authUser.FirstName,
			LastName:  authUser.LastName,
			Email:     authUser.Email,
			AvatarURL: authUser.AvatarURL,
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	token, err := s.TokenService.CreateToken(c, auth.ID, clientID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("created token: %s, '%s'", token.ID, token.PlainText)

	if err = sessionSetValue(c, TokenSessionKey, token.PlainText); err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	return c.Redirect(http.StatusFound, "/")
}

func (s *Server) AuthnValidator(tokenString string, c echo.Context) (bool, error) {
	token, err := s.getTokenFromSession(c)
	if token.ExpiresAt.Before(time.Now()) {
		log.Printf("token expired at %s\n", token.ExpiresAt)

		// bypass the transaction so the middleware doesn't roll back the token delete
		c.Set(keygo.ContextKeyTx, db.DB)
		if err = s.TokenService.DeleteToken(c, token.ID); err != nil {
			return false, fmt.Errorf("failed to delete expired token %s, %s", token.ID, err)
		}
		return false, nil
	}

	c.Set(keygo.ContextKeyToken, token)
	return true, nil
}

func AuthnSkipper(c echo.Context) bool {
	var skipURLs = []string{"/api/auth", "/api/auth/login", "/api/auth/callback", "/api/auth/logout"}
	for _, u := range skipURLs {
		if c.Path() == u {
			return true
		}
	}

	if c.Request().Method == http.MethodOptions {
		return true
	}
	return false
}

func (s *Server) getTokenFromSession(c echo.Context) (keygo.Token, error) {
	tokenInterface, err := sessionGetValue(c, TokenSessionKey)
	if err != nil {
		s.Logger.Infof("no token in session: %w\n", err)
		return keygo.Token{}, nil
	}

	tokenPlainText, ok := tokenInterface.(string)
	if !ok {
		return keygo.Token{}, fmt.Errorf("token in session is not a string\n")
	}

	token, err := s.TokenService.FindToken(c, tokenPlainText)
	if err != nil {
		return keygo.Token{}, fmt.Errorf("could not find token '%s' in DB: %w\n", tokenPlainText, err)
	}

	return token, nil
}
