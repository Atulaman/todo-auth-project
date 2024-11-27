package dbhelper

import (
	"database/sql"
	"time"
	"todo-auth/database"

	_ "github.com/lib/pq"
)

func CreateUser(username string, password string) error {

	_, err := database.TODO.Exec("INSERT INTO auth (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		return err
	}
	return nil
}

func IsUserExists(username, password string) error {
	var user string
	err := database.TODO.QueryRow("SELECT username FROM auth WHERE username = $1 AND password = $2", username, password).Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	return nil
}
func SetSession(username string, sessionID string) error {
	_, err := database.TODO.Exec("INSERT INTO session (session_id, username, created_at) VALUES ($1, $2, $3)", sessionID, username, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}
func DeleteSession(cookie string) error {
	_, err := database.TODO.Exec("DELETE FROM session WHERE session_id = $1", cookie)
	if err != nil {
		return err
	}
	return nil
}
