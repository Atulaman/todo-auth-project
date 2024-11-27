package middlewares

import (
	"database/sql"
	"net/http"
	"time"
	"todo-auth/database"
	"todo-auth/utils"

	_ "github.com/lib/pq"
)

func Caller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//cookie, err := r.Cookie("session_id")
		cookie, err := utils.GetSessionID(r)
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized user", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
			return
		}
		var username string
		var created_at time.Time
		err = database.TODO.QueryRow("SELECT username, created_at FROM session WHERE session_id = $1", cookie).Scan(&username, &created_at)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized user", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		duration := time.Now().UTC().Sub(created_at) //time.Since(created_at)
		if duration >= 5*time.Minute {
			_, err = database.TODO.Exec("DELETE FROM session WHERE session_id = $1", cookie)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Expires:  time.Unix(0, 0), // Expire the cookie immediately
				MaxAge:   -1,              // Set MaxAge to -1 to delete the cookie
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
