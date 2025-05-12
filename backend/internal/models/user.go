package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	PasswordHash string    `json:"-" db:"password_hash"`
	AuthProvider string    `json:"auth_provider" db:"auth_provider"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`

	// Fields for JSON marshaling
	CreatedAtEpoch int64 `json:"created_at" db:"-"`
	UpdatedAtEpoch int64 `json:"updated_at" db:"-"`
}

// PrepareMarshal prepares the User for JSON marshaling by setting epoch timestamps
func (u *User) PrepareMarshal() {
	// Convert time.Time to Unix timestamps in milliseconds
	u.CreatedAtEpoch = u.CreatedAt.UnixNano() / int64(time.Millisecond)
	u.UpdatedAtEpoch = u.UpdatedAt.UnixNano() / int64(time.Millisecond)
}

// MarshalJSON implements a custom JSON marshaler for User
func (u *User) MarshalJSON() ([]byte, error) {
	// Prepare the timestamps if they haven't been set
	if u.CreatedAtEpoch == 0 && !u.CreatedAt.IsZero() {
		u.CreatedAtEpoch = u.CreatedAt.UnixNano() / int64(time.Millisecond)
	}

	if u.UpdatedAtEpoch == 0 && !u.UpdatedAt.IsZero() {
		u.UpdatedAtEpoch = u.UpdatedAt.UnixNano() / int64(time.Millisecond)
	}

	// Create an alias to avoid infinite recursion
	type Alias User
	return json.Marshal(&struct {
		*Alias
		CreatedAt int64 `json:"created_at"`
		UpdatedAt int64 `json:"updated_at"`
	}{
		Alias:     (*Alias)(u),
		CreatedAt: u.CreatedAtEpoch,
		UpdatedAt: u.UpdatedAtEpoch,
	})
}
