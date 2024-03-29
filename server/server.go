package server

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server struct {
	*echo.Echo
	db *gorm.DB
}

const loggerFormat = "${time_rfc3339} ${status} ${method} ${uri} ${error}\n"

var svr *Server

type Option func(*Server)

func WithDataBase(db *gorm.DB) Option {
	return func(s *Server) {
		s.db = db
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

	// Logger Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: loggerFormat}))

	// Recover Middleware
	e.Use(middleware.Recover())

	// Session Middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	// DB Transaction Middleware
	e.Use(TxMiddleware(svr.db))

	if os.Getenv("GO_ENV") == "development" {
		e.Debug = true
	}

	// serve static assets, e.g. favicon.ico
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Browse: false,
	}))

	svr.registerAPIRoutes()

	// send all other routes to the UI router
	svr.registerUIRoutes()

	return svr
}

func (s *Server) registerAPIRoutes() {
	api := s.Group("/api", s.AuthnMiddleware)

	api.GET("/auth", s.authStatus)
	api.GET("/auth/login", s.authLogin)
	api.GET("/auth/callback", s.authCallback)
	api.GET("/auth/logout", s.authLogout)

	api.POST("/tenants", s.tenantsCreateHandler)
	api.GET("/tenants", s.tenantsListHandler)
	api.GET("/tenants/:id", s.tenantsGetHandler)

	api.POST("/tenants/:id/users", s.tenantsUsersCreateHandler)

	api.GET("/users", s.usersListHandler)
	api.GET("/users/:id", s.userHandler)
	api.PUT("/users/:id", s.usersUpdateHandler)
}
