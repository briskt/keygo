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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${uri} ${error}\n",
	}))

	// Recover Middleware
	e.Use(middleware.Recover())

	// Session Middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	// DB Transaction Middleware
	e.Use(TxMiddleware(db.DB))

	if os.Getenv("GO_ENV") == "development" {
		e.Debug = true
	}

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Browse: false,
	}))

	svr.registerRoutes()
	svr.getServices()
	return svr
}

func (s *Server) registerRoutes() {
	api := s.Group("/api", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:   AuthnSkipper,
		Validator: s.AuthnValidator,
	}))
	api.GET("/auth", s.authStatus)
	api.GET("/auth/login", s.authLogin)
	api.GET("/auth/callback", s.authCallback)
	api.GET("/auth/logout", s.authLogout)
	api.GET("/users/:id", s.userHandler)

	s.registerUiRoutes()
}

func (s *Server) getServices() {
	s.UserService = db.NewUserService()
	s.AuthService = db.NewAuthService()
	s.TokenService = db.NewTokenService()
}
