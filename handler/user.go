package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
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
	err := utils.DecodeJson(r, &user)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		log.Logging(err, "Error decoding request body", 400, "warning", r)
		return
	}
	if user.Username == "" || user.Password == "" || utf8.RuneCountInString(user.Username) > 20 || utf8.RuneCountInString(user.Password) > 20 || utf8.RuneCountInString(user.Password) < 8 || utf8.RuneCountInString(user.Username) < 8 {
		utils.ResponseError(w, "Missing/Invalid username or password", http.StatusBadRequest)
		log.Logging(err, "Missing/Invalid username or password", 400, "warning", r)
		return
	}
	err = dbhelper.CreateUser(user.Username, user.Password)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error creating user", 500, "error", r)
		return
	}
	utils.ResponseJson(w, http.StatusOK, map[string]interface{}{"message": "Registration successful"})
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
	err := utils.DecodeJson(r, &user)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusBadRequest)
		log.Logging(err, "Error decoding request body", 400, "warning", r)
		return
	}
	if user.Username == "" || user.Password == "" {
		utils.ResponseError(w, "Missing username or password", http.StatusBadRequest)
		log.Logging(err, "Missing username or password", 400, "warning", r)
		return
	}
	err = dbhelper.IsUserExists(user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseError(w, "Invalid username or password", http.StatusUnauthorized)
			log.Logging(err, "Invalid username or password", 401, "warning", r)
			return
		}
		utils.ResponseError(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error checking user", 500, "error", r)
		return
	}
	sessionID, err := generateSessionID()
	if err != nil {
		utils.ResponseError(w, "Error generating session ID", http.StatusInternalServerError)
		log.Logging(err, "Error generating session ID", 500, "error", r)
		return
	}
	err = dbhelper.SetSession(user.Username, sessionID)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusInternalServerError)
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
	utils.ResponseJson(w, http.StatusOK, map[string]interface{}{"message": "Login successful"})
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
	cookie, err := utils.GetSessionID(r)
	if err != nil {
		if err == http.ErrNoCookie {
			utils.ResponseError(w, "Already logged out", http.StatusUnauthorized)
			log.Logging(err, "Already logged out", 401, "warning", r)
			return
		}
		utils.ResponseError(w, "Error retrieving cookie", http.StatusInternalServerError)
		log.Logging(err, "Error retrieving cookie", http.StatusInternalServerError, "error", r)
		return
	}

	err = dbhelper.DeleteSession(cookie)
	if err != nil {
		utils.ResponseError(w, "Error deleting session", http.StatusInternalServerError)
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

	utils.ResponseJson(w, http.StatusOK, map[string]interface{}{"message": "Logout successful"})
	log.Logging(nil, "Logout successful", 200, "info", r)
}
