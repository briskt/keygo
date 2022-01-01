package http

import (
	"fmt"
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

func init() {
	// Goth
	goth.UseProviders(
		auth0.New(os.Getenv("AUTH0_KEY"),
			os.Getenv("AUTH0_SECRET"),
			os.Getenv("HOST")+AuthCallbackPath+"?provider=auth0",
			os.Getenv("AUTH0_DOMAIN"),
		),
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			os.Getenv("HOST")+AuthCallbackPath+"?provider=google",
		),
	)
}

func RegisterAuthRoutes(e *echo.Echo) {
	// Route => handler
	e.POST("/auth/login", authLogin)
	e.GET(AuthCallbackPath, authCallback)
	e.GET("/auth/logout", authLogout)
}

func authLogin(c echo.Context) error {
	res := c.Response()
	req := c.Request()

	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		return c.JSON(http.StatusOK, gothUser)
	}

	fmt.Printf("client ID: %s\n", c.QueryParam("client_id"))
	allProviders := map[string]string{"auth0": "Auth0", "google": "Google"}
	options := make([]ProviderOption, len(allProviders))
	i := 0
	for key, name := range allProviders {
		options[i] = ProviderOption{Key: key, Name: name}

		gothic.GetProviderName = func(req *http.Request) (string, error) {
			return key, nil
		}
		url, err := gothic.GetAuthURL(res, req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, AuthError{Error: err.Error()})
		}
		options[i].RedirectURL = url
		i++
	}
	return c.JSON(http.StatusOK, options)
}

func authLogout(c echo.Context) error {
	gothic.Logout(c.Response(), c.Request())
	c.Response().Header().Set("Location", "/")
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func authCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		_, err = fmt.Fprintln(c.Response(), err)
		return err
	}
	return c.JSON(http.StatusOK, user)
}
