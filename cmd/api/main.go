package main

import (
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
	"github.com/briskt/keygo/server"
)

func main() {
	fmt.Println("starting API")

	dbConnection := db.OpenDB()
	services := app.DataServices{
		AuthService:   db.NewAuthService(),
		TenantService: db.NewTenantService(),
		TokenService:  db.NewTokenService(),
		UserService:   db.NewUserService(),
	}
	e := server.New(server.WithDataBase(dbConnection), server.WithDataServices(services))

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
