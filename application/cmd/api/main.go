package main

import (
	"fmt"

	"github.com/schparky/keygo/server"
)

func main() {
	fmt.Println("starting API")

	e := server.New()

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
