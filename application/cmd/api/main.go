package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/schparky/keygo/http"
)

func main() {
	fmt.Println("starting API")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Gorm DB Middleware
	dsn := "host=localhost user=keygo password=keygo dbname=keygo port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error opening database, " + err.Error())
	}
	e.Use(http.Transaction(db))

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{os.Getenv("UI_URL")},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
	}))

	http.RegisterAuthRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
