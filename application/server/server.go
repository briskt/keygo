package server

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
)

type Server struct {
	*echo.Echo

	TokenService keygo.TokenService
	AuthService  keygo.AuthService
	UserService  keygo.UserService
}

var svr *Server

func New() *Server {
	if svr != nil {
		return svr
	}
	e := echo.New()
	svr = &Server{Echo: e}

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
		Validator: svr.AuthnMiddleware,
	}))

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{os.Getenv("UI_URL")},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
	}))

	if os.Getenv("GO_ENV") == "development" {
		e.Debug = true
	}

	svr.registerRoutes()
	svr.getServices()
	return svr
}

func (s *Server) registerRoutes() {
	s.registerAuthRoutes()
	// s.registerFormRoutes()
	s.registerUserRoutes()
}

func (s *Server) getServices() {
	s.UserService = db.NewUserService()
	s.AuthService = db.NewAuthService()
	s.TokenService = db.NewTokenService()
}
