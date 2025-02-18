package dbhelper

import (
	"database/sql"
	"time"
	"todo-auth/database"

	_ "github.com/lib/pq"
)

func CreateUser(username string, password string) error {
	query := `INSERT INTO auth (username, password) VALUES ($1, $2)`
	_, err := database.TODO.Exec(query, username, password)
	if err != nil {
		return err
	}
	return nil
}

func IsUserExists(username, password string) error {
	var user string
	query := `SELECT username FROM auth WHERE username = $1 AND password = $2`
	err := database.TODO.Get(&user, query, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	return nil
}
func SetSession(username string, sessionID string) error {
	query := `INSERT INTO session (session_id, username, created_at) VALUES ($1, $2, $3)`
	_, err := database.TODO.Exec(query, sessionID, username, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}
func DeleteSession(cookie string) error {
	query := `DELETE FROM session WHERE session_id = $1`
	_, err := database.TODO.Exec(query, cookie)
	if err != nil {
		return err
	}
	return nil
}
