package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (s *Server) usersListHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != app.UserRoleAdmin {
		return echo.NewHTTPError(http.StatusOK, []app.User{})
	}

	users, err := db.FindUsers(c, app.UserFilter{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("found %d users", len(users))

	return c.JSON(http.StatusOK, users)
}

func (s *Server) userHandler(c echo.Context) error {
	actor := app.CurrentUser(c)
	var user app.User

	id := c.Param("id")
	if id == actor.ID {
		user = actor
	} else if actor.Role != app.UserRoleAdmin {
		return echo.NewHTTPError(http.StatusNotFound, AuthError{Error: "not found"})
	} else {
		dbUser, err := db.FindUserByID(c, id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		user, err = db.ConvertUser(c, dbUser)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) usersUpdateHandler(c echo.Context) error {
	var input app.UserUpdateInput
	err := (&echo.DefaultBinder{}).BindBody(c, &input)
	if err != nil {
		// TODO: improve error response here (and probably everywhere else too)
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	id := c.Param("id")
	actor := app.CurrentUser(c)
	if id != actor.ID && actor.Role != app.UserRoleAdmin {
		return echo.NewHTTPError(http.StatusNotFound, AuthError{Error: "not found"})
	}

	updatedUser, err := db.UpdateUser(c, id, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user, err := db.ConvertUser(c, updatedUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, user)
}
