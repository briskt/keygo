package main

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/schparky/keygo/db"
	"github.com/schparky/keygo/http"
)

func main() {
	fmt.Println("starting API")

	// Echo instance
	e := echo.New()

	// Logger Middleware
	e.Use(middleware.Logger())

	// Recover Middleware
	e.Use(middleware.Recover())

	// Session Middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	// Gorm DB Middleware
	e.Use(http.Transaction(db.DB))

	// Authn Middleware
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:   http.AuthnSkipper,
		Validator: http.AuthnMiddleware,
	}))

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{os.Getenv("UI_URL")},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
	}))

	http.RegisterAuthRoutes(e)
	http.RegisterUserRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
