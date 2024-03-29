package oauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
	scopes string
}

// Manager is an interface that defines a user identity manager service
// TODO: implement this for Authgear
type Manager interface {
	// ReadUser requests the user profile from the ID Manager database
	ReadUser(userID string) (Profile, error)

	// ResendEmailVerificationMessage requests a new copy of the email verification message for the user's currently
	// registered email address.
	ResendEmailVerificationMessage(userID, email string) error

	// UpdateEmail update's the user's email address in the ID Manager database. The address should be verified using
	// RequestEmailVerificationCode and ValidateCode prior to changing to a new address.
	UpdateEmail(userID, email string) error

	// DeleteUser deletes a user from the auth manager's database
	DeleteUser(userID string) error
}

type Profile struct {
	ID       string
	Email    string
	Verified bool
}

type Config struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       string
}

var authenticator *Authenticator

// Init initializes the Authenticator.
func Init(config Config) error {
	if authenticator != nil {
		return nil
	}

	// initialize the provider using service discovery
	provider, err := oidc.NewProvider(
		context.Background(),
		config.Issuer,
	)
	if err != nil {
		return fmt.Errorf("oauth initialization error: %w", err)
	}

	conf := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	authenticator = &Authenticator{
		Provider: provider,
		Config:   conf,
		scopes:   config.Scopes,
	}

	return nil
}

func Get() *Authenticator {
	return authenticator
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

// AuthCodeURL is a wrapper for oauth2.AuthCodeURL, which returns a URL to OAuth 2.0 provider's consent page
// that asks for permissions for the required scopes explicitly.
func AuthCodeURL(state string) string {
	options := oauth2.SetAuthURLParam("scope", authenticator.scopes)

	return authenticator.AuthCodeURL(state, options)
}

func GetProfile(ctx context.Context, code string) (Profile, error) {
	var ap Profile

	// Exchange an authorization code for a token.
	token, err := authenticator.Exchange(ctx, code)
	if err != nil {
		err = fmt.Errorf("failed to create an access token from the authorization code: %w", err)
		return ap, err
	}

	idToken, err := authenticator.VerifyIDToken(ctx, token)
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
