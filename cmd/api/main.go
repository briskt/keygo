package main

import (
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/briskt/keygo/db"
	"github.com/briskt/keygo/server"
)

func main() {
	fmt.Println("starting API")

	dbConnection := db.OpenDB()
	e := server.New(server.WithDataBase(dbConnection))

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
