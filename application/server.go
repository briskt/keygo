package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/google"
	"html/template"
	"net/http"
	"os"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/login?provider={{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Goth
	goth.UseProviders(
		auth0.New(os.Getenv("AUTH0_KEY"), os.Getenv("AUTH0_SECRET"), "http://localhost:3000/auth/auth0/callback", os.Getenv("AUTH0_DOMAIN")),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
	)

	// Route => handler
	e.GET("/", home)
	e.GET("/auth/login", authLogin)
	e.GET("/auth/callback", authCallback)
	e.GET("/logout", authLogout)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func home(c echo.Context) error {
	m := map[string]string{"auth0": "Auth0", "google": "Google"}
	keys := []string{"auth0", "google"}
	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(c.Response(), providerIndex)
	return nil
}

func authLogin(c echo.Context) error {
	//try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(c.Response(), gothUser)
	} else {
		gothic.BeginAuthHandler(c.Response(), c.Request())
	}

	return nil
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
	t, _ := template.New("foo").Parse(userTemplate)
	return t.Execute(c.Response(), user)
}
