package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var uiRoutes = []string{"/"}

func (s *Server) registerUIRoutes() {
	s.GET("*", spaHandler)
}

func spaHandler(c echo.Context) error {
	fs := http.Dir("public")
	file, err := fs.Open("./index.html")
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	http.ServeContent(c.Response(), c.Request(), info.Name(), info.ModTime(), file)
	return nil
}
