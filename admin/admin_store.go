package admin

import (
	"database/sql"

	"github.com/Sairenu1/movie_booking/store"
)

type Admin struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate admin login
func ValidateAdmin(username, password string) bool {
	var storedPass string

	err := store.DB.QueryRow(`
		SELECT password FROM admins WHERE username = ?`,
		username,
	).Scan(&storedPass)

	if err == sql.ErrNoRows {
		return false
	}

	if err != nil {
		return false
	}

	// Plain password check (simple version)
	return storedPass == password
}
