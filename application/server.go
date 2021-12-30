package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/google"
	"net/http"
	"os"
)

const AuthCallback = "/auth/callback"

type AuthError struct {
	Error string
}

type ProviderOption struct {
	Key         string
	Name        string
	RedirectURL string
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{os.Getenv("UI_URL")},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
	}))

	// Goth
	goth.UseProviders(
		auth0.New(os.Getenv("AUTH0_KEY"),
			os.Getenv("AUTH0_SECRET"),
			os.Getenv("HOST")+AuthCallback+"?provider=auth0",
			os.Getenv("AUTH0_DOMAIN"),
		),
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			os.Getenv("HOST")+AuthCallback+"?provider=google",
		),
	)

	// Route => handler
	e.POST("/auth/login", authLogin)
	e.GET(AuthCallback, authCallback)
	e.GET("/auth/logout", authLogout)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
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
