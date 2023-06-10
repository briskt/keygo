package main

import (
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/briskt/keygo/server"
)

func main() {
	fmt.Println("starting API")

	e := server.New()

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
