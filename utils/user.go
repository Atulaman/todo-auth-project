package utils

import (
	"encoding/json"
	"net/http"
)

func GetSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", http.ErrNoCookie
		}
		return "", err
	}
	return cookie.Value, err
}
func DecodeJson(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
func ResponseJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
func ResponseError(w http.ResponseWriter, errmsg string, code int) {
	http.Error(w, errmsg, code)
}
