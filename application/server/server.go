package server

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/schparky/keygo/db"
)

type Server struct {
	*echo.Echo
}

var svr *Server

func New() *Server {
	if svr != nil {
		return svr
	}

	// Echo instance
	e := echo.New()

	// Logger Middleware
	e.Use(middleware.Logger())

	// Recover Middleware
	e.Use(middleware.Recover())

	// Session Middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	// DB Transaction Middleware
	e.Use(TxMiddleware(db.DB))

	// Authn Middleware
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:   AuthnSkipper,
		Validator: AuthnMiddleware,
	}))

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{os.Getenv("UI_URL")},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
	}))

	RegisterAuthRoutes(e)
	RegisterFormRoutes(e)
	RegisterUserRoutes(e)

	svr = &Server{Echo: e}
	return svr
}
