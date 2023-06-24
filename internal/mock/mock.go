package mock

import (
	"strconv"
)

var newID func() string

var lastID = 0

func init() {
	newID = func() string {
		lastID++
		return strconv.Itoa(lastID)
	}
}
