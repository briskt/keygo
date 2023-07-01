package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/server/oauth"
)

const AuthCallbackPath = "/api/auth/callback"

const (
	SessionKeyAuthID   = "AuthID"
	SessionKeyToken    = "Token"
	SessionKeyReturnTo = "ReturnTo"
	SessionKeyState    = "State"
)

const (
	ParamReturnTo = "returnTo"
	ParamCode     = "code"
	DefaultUIPath = "/"
)

type AuthError struct {
	Error string
}

func init() {
	const required = true
	config := oauth.Config{
		Issuer:       env("OAUTH_ISSUER_URL", required),
		ClientID:     env("OAUTH_CLIENT_ID", required),
		ClientSecret: env("OAUTH_CLIENT_SECRET", required),
		RedirectURL:  env("HOST", required) + env("OAUTH_REDIRECT_PATH", required),
		Scopes:       env("OAUTH_OPENID_SCOPES", required),
	}

	if err := oauth.Init(config); err != nil {
		log.Fatalf("error initializing authenticator: %s", err)
	}
}

func (s *Server) authStatus(c echo.Context) error {
	var status app.AuthStatus

	token, err := s.getTokenFromSession(c)
	if err != nil {
		s.Logger.Warnf("error getting token from session: %s", err)
		return c.JSON(http.StatusOK, status)
	}

	if token.ExpiresAt.After(time.Now()) {
		status.IsAuthenticated = true
	}
	status.UserID = token.UserID
	status.Expiry = token.ExpiresAt

	return c.JSON(http.StatusOK, status)
}

func (s *Server) authLogin(c echo.Context) error {
	authenticator := oauth.Get()
	if authenticator == nil {
		err := fmt.Errorf("authenticator is not initialized")
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	returnToPath, err := getReturnTo(c)
	if err != nil {
		s.Logger.Warnf("getting return path: %w", err)
		returnToPath = DefaultUIPath
	}

	if err := sessionSetValue(c, SessionKeyReturnTo, returnToPath); err != nil {
		err = fmt.Errorf("setting return path: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	state := generateRandomState()
	if err := sessionSetValue(c, SessionKeyState, state); err != nil {
		err = fmt.Errorf("generating state key: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	url := oauth.AuthCodeURL(state)
	s.Logger.Infof("redirecting to auth provider: %s", url)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func getReturnTo(c echo.Context) (string, error) {
	if path := c.Param(ParamReturnTo); path != "" {
		return path, nil
	}

	path, err := sessionGetString(c, SessionKeyReturnTo)
	if err != nil {
		return "", fmt.Errorf("failed to get %s from session: %w", SessionKeyReturnTo, err)
	}
	return path, nil
}

func generateRandomState() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate session state, rand.Read returned error: " + err.Error())
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state
}

func (s *Server) authLogout(c echo.Context) error {
	token, err := s.getTokenFromSession(c)
	if err != nil {
		s.Logger.Error(err.Error())
	}
	if err := s.TokenService.DeleteToken(c, token.ID); err != nil {
		s.Logger.Errorf("failed to delete user token: %s", err)
	}
	return c.Redirect(http.StatusTemporaryRedirect, DefaultUIPath)
}

func (s *Server) authCallback(c echo.Context) error {
	if authError := c.QueryParam("error"); authError != "" {
		errDescription := c.QueryParam("error_description")
		err := fmt.Errorf("auth error: %s, description: %s", authError, errDescription)
		s.Logger.Errorf("auth error %s: %s", authError, errDescription)
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	profile, err := getAuthProfile(c)
	if err != nil {
		err = fmt.Errorf("auth profile error: %w", err)
		s.Logger.Errorf(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	if err := sessionSetValue(c, SessionKeyAuthID, profile.ID); err != nil {
		err = fmt.Errorf("setting authID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("user authenticated, profile=%+v", profile)

	token, err := s.TokenService.CreateToken(c, app.Token{
		AuthID: profile.ID,
		User: app.User{
			Email: profile.Email,
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("created token: %s", token.ID)

	if err := s.UserService.TouchLastLoginAt(c, token.UserID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	if err = sessionSetValue(c, SessionKeyToken, token.PlainText); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	returnTo, err := getReturnTo(c)
	if err != nil {
		err = fmt.Errorf("getting return path: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	return c.Redirect(http.StatusTemporaryRedirect, returnTo)
}

func (s *Server) AuthnMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if AuthnSkipper(c) {
			return next(c)
		}

		token, err := s.getTokenFromSession(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("getTokenFromSession: %s", err))
		}

		if token.ID == "" {
			token, _ = s.TokenService.FindToken(c, getBearerToken(c))
		}

		status := http.StatusUnauthorized
		authError := AuthError{"not authorized"}
		if token.ID == "" {
			return echo.NewHTTPError(status, authError)
		}
		if token.ExpiresAt.Before(time.Now()) {
			s.Logger.Infof("token expired at %s\n", token.ExpiresAt)
			return echo.NewHTTPError(status, authError)
		}

		if err := s.TokenService.UpdateToken(c, token.ID); err != nil {
			return echo.NewHTTPError(status, echo.NewHTTPError(http.StatusInternalServerError), AuthError{Error: err.Error()})
		}

		c.Set(app.ContextKeyToken, token)
		c.Set(app.ContextKeyUser, token.User)
		return next(c)
	}
}

func getBearerToken(c echo.Context) (token string) {
	for _, h := range c.Request().Header["Authorization"] {
		parts := strings.Split(h, " ")
		if len(parts) != 2 {
			continue
		}
		if strings.ToLower(parts[0]) == "bearer" {
			token = parts[1]
		}
	}
	return token
}

func AuthnSkipper(c echo.Context) bool {
	skipURLs := []string{"/api/auth", "/api/auth/login", "/api/auth/callback", "/api/auth/logout"}
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

func (s *Server) getTokenFromSession(c echo.Context) (app.Token, error) {
	tokenInterface, err := sessionGetValue(c, SessionKeyToken)
	if err != nil {
		s.Logger.Infof("no token in session: %s", err)
		return app.Token{}, nil
	}

	tokenPlainText, ok := tokenInterface.(string)
	if !ok {
		return app.Token{}, fmt.Errorf("token in session is not a string\n")
	}

	token, err := s.TokenService.FindToken(c, tokenPlainText)
	if err != nil {
		return app.Token{}, fmt.Errorf("could not find token in DB: %w\n", err)
	}

	return token, nil
}

// getAuthProfile retrieves the user's profile from the authorization code provided by Auth0. This code is based on the
// Auth0 Quick Start at https://auth0.com/docs/quickstart/webapp/golang/interactive
func getAuthProfile(c echo.Context) (oauth.Profile, error) {
	var ap oauth.Profile

	// TODO: move this whole function to the oauth package?

	authenticator := oauth.Get()
	if authenticator == nil {
		err := fmt.Errorf("authenticator is not initialized")
		return ap, err
	}

	// Exchange an authorization code for a token.
	token, err := authenticator.Exchange(context.TODO(), c.QueryParam(ParamCode)) // TODO: need a real Context?
	if err != nil {
		err = fmt.Errorf("failed to create an access token from the authorization code: %w", err)
		return ap, err
	}

	idToken, err := authenticator.VerifyIDToken(context.TODO(), token) // TODO: need a real context?
	if err != nil {
		err = fmt.Errorf("failed to verify ID Token: %w", err)
		return ap, err
	}

	var profile map[string]any
	if err = idToken.Claims(&profile); err != nil {
		err = fmt.Errorf("failed to get profile from token: %w", err)
		return ap, err
	}

	var ok bool

	ap.ID, ok = profile["sub"].(string)
	if !ok {
		err = fmt.Errorf("no 'sub' key (AuthID) found in the user profile")
		return ap, err
	}

	ap.Email, ok = profile["email"].(string)
	if !ok {
		err = fmt.Errorf("no email address found in the user profile")
		return ap, err
	}

	ap.Verified, ok = profile["email_verified"].(bool)
	if !ok {
		err = fmt.Errorf("invalid email_verified in the user profile")
		return ap, err
	}

	return ap, nil
}

func env(key string, required bool) string {
	v := os.Getenv(key)
	if v == "" && required {
		panic("required environment variable '" + key + "' is not defined")
	}
	return v
}
