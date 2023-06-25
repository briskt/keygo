package app

import (
	"time"
)

// swagger:model
type AuthStatus struct {
	// IsAuthenticated is true when the supplied session cookie is valid and references a valid user
	IsAuthenticated bool `json:"IsAuthenticated"`

	// Expiry is the date and time when the session is scheduled to expire. It is invalid if `IsAuthenticated` is false.
	//
	// swagger:strfmt date-time
	Expiry time.Time `json:"Expiry"`

	// UserID is the ID of the authenticated user. It is invalid if `IsAuthenticated` is false.
	UserID string `json:"UserID"`
}
