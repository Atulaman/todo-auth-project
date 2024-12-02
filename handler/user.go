package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
	dbhelper "todo-auth/database/db-helper"
	log "todo-auth/logging"
	"todo-auth/utils"
	"unicode/utf8"

	_ "github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User registration data"
// @Success 200 {object} map[string]interface{} "Registration successful"
// @Failure 400 {object} map[string]string "Invalid username or password"
// @Failure 500 {object} map[string]string "Error inserting user or user already exists"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Logging(err, "Error decoding request body", 400, "warning", r)
		return
	}
	if user.Username == "" || user.Password == "" || utf8.RuneCountInString(user.Username) > 20 || utf8.RuneCountInString(user.Password) > 20 || utf8.RuneCountInString(user.Password) < 8 || utf8.RuneCountInString(user.Username) < 8 {
		http.Error(w, "Missing/Invalid username or password", http.StatusBadRequest)
		log.Logging(err, "Missing/Invalid username or password", 400, "warning", r)
		return
	}
	//_, err = database.TODO.Exec("INSERT INTO auth (username, password) VALUES ($1, $2)", user.Username, user.Password)
	err = dbhelper.CreateUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error creating user", 500, "error", r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Registration successful"})
	log.Logging(nil, "Registration successful", 200, "info", r)
}
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and create a session
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]string "Invalid username or password"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Error logging in"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Logging(err, "Error decoding request body", 400, "warning", r)
		return
	}
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		log.Logging(err, "Missing username or password", 400, "warning", r)
		return
	}
	// var (
	// 	username string
	// 	password string
	// )
	//err = database.TODO.QueryRow("SELECT username,password FROM auth WHERE username = $1 AND password = $2", user.Username, user.Password).Scan(&username, &password)
	err = dbhelper.IsUserExists(user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			log.Logging(err, "Invalid username or password", 401, "warning", r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error checking user", 500, "error", r)
		return
	}
	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		log.Logging(err, "Error generating session ID", 500, "error", r)
		return
	}
	//_, err = database.TODO.Exec("INSERT INTO session (session_id, username, created_at) VALUES ($1, $2, $3)", sessionID, user.Username, time.Now().UTC())
	err = dbhelper.SetSession(user.Username, sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error creating session", 500, "error", r)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	// if r.URL.Path == "/login" {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Login successful"})
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Login successful"})
	log.Logging(nil, "Login successful", 200, "info", r)
}

// Logout godoc
// @Summary Logout a user
// @Description Logout a user by invalidating their session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]string "Already logged out or invalid session"
// @Failure 500 {object} map[string]string "Error deleting session"
// @Router /logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	//cookie, err := r.Cookie("session_id")
	cookie, err := utils.GetSessionID(r)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Already logged out", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
		return
	}

	//_, err = database.TODO.Exec("DELETE FROM session WHERE session_id = $1", cookie)
	err = dbhelper.DeleteSession(cookie)
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
	if r.URL.Path == "/logout" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Logout successful"})
	}
	log.Logging(nil, "Logout successful", 200, "info", r)
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{"message": "Logout successful"})
}
