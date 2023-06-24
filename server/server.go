package server

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

type Server struct {
	*echo.Echo
	mode string

	AuthService   app.AuthService
	TenantService app.TenantService
	TokenService  app.TokenService
	UserService   app.UserService
}

const loggerFormat = "${time_rfc3339} ${status} ${method} ${uri} ${error}\n"

const test = "test"

var svr *Server

type Option func(*Server)

func TestMode() Option {
	return func(s *Server) {
		s.mode = test
	}
}

func New(options ...Option) *Server {
	if svr != nil {
		return svr
	}
	e := echo.New()
	svr = &Server{Echo: e}

	for _, opt := range options {
		opt(svr)
	}

	if svr.mode != test {
		svr.getServices()
	}

	// Logger Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: loggerFormat}))

	// Recover Middleware
	e.Use(middleware.Recover())

	// Session Middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	if svr.mode != test {
		// DB Transaction Middleware
		e.Use(TxMiddleware(db.DB))
	}

	if os.Getenv("GO_ENV") == "development" {
		e.Debug = true
	}

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Browse: false,
	}))

	svr.registerRoutes()
	return svr
}

func (s *Server) registerRoutes() {
	api := s.Group("/api", s.AuthnMiddleware)

	api.GET("/auth", s.authStatus)
	api.GET("/auth/login", s.authLogin)
	api.GET("/auth/callback", s.authCallback)
	api.GET("/auth/logout", s.authLogout)
	api.POST("/tenants", s.tenantsCreateHandler)
	api.GET("/tenants", s.tenantsListHandler)
	api.GET("/users", s.usersListHandler)
	api.GET("/users/:id", s.userHandler)

	s.registerUiRoutes()
}

func (s *Server) getServices() {
	s.AuthService = db.NewAuthService()
	s.TenantService = db.NewTenantService()
	s.TokenService = db.NewTokenService()
	s.UserService = db.NewUserService()
}
