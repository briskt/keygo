package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/server/oauth"
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

// TODO: make a connection between error responses and the logger, so all errors get logged
type AuthError struct {
	Error string
}

func init() {
	config := oauth.Config{
		Issuer:       os.Getenv("OAUTH_ISSUER_URL"),
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		Scopes:       os.Getenv("OAUTH_OPENID_SCOPES"),
	}

	// TODO: check here for missing config

	if err := oauth.Init(config); err != nil {
		log.Fatalf("error initializing authenticator: %s", err)
	}
}

func (s *Server) authStatus(c echo.Context) error {
	var status keygo.AuthStatus

	token, err := s.getTokenFromSession(c)
	if err != nil {
		s.Logger.Warnf("error getting token from session: %s", err)
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
	authenticator := oauth.Get()
	if authenticator == nil {
		err := fmt.Errorf("authenticator is not initialized")
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	returnToPath, err := getReturnTo(c)
	if err != nil {
		s.Logger.Warnf("getting return path: %w", err)
		returnToPath = DefaultUIPath
	}

	if err := sessionSetValue(c, SessionKeyReturnTo, returnToPath); err != nil {
		err = fmt.Errorf("setting return path: %w", err)
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	state := generateRandomState()
	if err := sessionSetValue(c, SessionKeyState, state); err != nil {
		err = fmt.Errorf("generating state key: %w", err)
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
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
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	profile, err := getAuthProfile(c)
	if err != nil {
		err = fmt.Errorf("auth profile error: %w", err)
		s.Logger.Errorf(err.Error())
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	if err := sessionSetValue(c, SessionKeyAuthID, profile.ID); err != nil {
		err = fmt.Errorf("setting authID: %w", err)
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("user authenticated, profile=%+v", profile)

	auth, err := s.AuthService.CreateAuth(c, keygo.Auth{
		Provider:   "oauth",
		ProviderID: profile.ID,
		User: keygo.User{
			Email: profile.Email,
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	token, err := s.TokenService.CreateToken(c, auth.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("created token: %s, '%s'", token.ID, token.PlainText)

	if err = sessionSetValue(c, SessionKeyToken, token.PlainText); err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	returnTo, err := getReturnTo(c)
	if err != nil {
		err = fmt.Errorf("getting return path: %w", err)
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
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
			s.Logger.Error(err.Error())
		}
		if token.ExpiresAt.Before(time.Now()) {
			log.Printf("token expired at %s\n", token.ExpiresAt)
			return c.JSON(http.StatusNotFound, AuthError{"not found"})
		}

		c.Set(keygo.ContextKeyToken, token)
		c.Set(keygo.ContextKeyUser, token.Auth.User)
		return next(c)
	}
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

func (s *Server) getTokenFromSession(c echo.Context) (keygo.Token, error) {
	tokenInterface, err := sessionGetValue(c, SessionKeyToken)
	if err != nil {
		s.Logger.Infof("no token in session: %s", err)
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

// getAuthProfile retrieves the user's profile from the authorization code provided by Auth0. This code is based on the
// Auth0 Quick Start at https://auth0.com/docs/quickstart/webapp/golang/interactive
func getAuthProfile(c echo.Context) (oauth.Profile, error) {
	var ap oauth.Profile

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
