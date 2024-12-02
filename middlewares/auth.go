package middlewares

import (
	"database/sql"
	"net/http"
	"time"
	"todo-auth/database"
	log "todo-auth/logging"
	"todo-auth/utils"

	_ "github.com/lib/pq"
)

func Caller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := utils.GetSessionID(r)
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized user", http.StatusUnauthorized)
				log.Logging(err, "Unauthorized user", 401, "warning", r)
				return
			}
			http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
			log.Logging(err, "Error retrieving cookie", 500, "error", r)
			return
		}
		//var username string
		var created_at time.Time
		data := &struct {
			Username  string    `db:"username"`
			CreatedAt time.Time `db:"created_at"`
		}{}
		query := `SELECT username, created_at
		FROM session
		WHERE session_id = $1`
		err = database.TODO.Get(data, query, cookie)
		created_at = data.CreatedAt
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized user", http.StatusUnauthorized)
				log.Logging(err, "Unauthorized user", 401, "warning", r)
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
				log.Logging(err, "Error deleting session", 500, "error", r)
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
			log.Logging(err, "Unauthorized user", 401, "warning", r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
