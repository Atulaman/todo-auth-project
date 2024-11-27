package utils

import "net/http"

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
