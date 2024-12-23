package middlewares

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func SetDB(DB *sql.DB) {
	db = DB
	//db = utils.GetDb()
}
func Caller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_id")
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
		err = db.QueryRow("SELECT username, created_at FROM session WHERE session_id = $1", cookie.Value).Scan(&username, &created_at)
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
			_, err = db.Exec("DELETE FROM session WHERE session_id = $1", cookie.Value)
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
