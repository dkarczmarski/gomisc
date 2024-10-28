package htserver

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type User struct {
	Username string
	Password string
}

var users []User = []User{
	{
		"admin", "123",
	},
}

func createAuthFunc(users []User) func(username, password string) *User {
	return func(username, password string) *User {
		for _, u := range users {
			if u.Username == username && u.Password == password {
				user := u
				return &user
			}
		}
		return nil
	}
}

func HandleBasicAuth(w http.ResponseWriter, r *http.Request, authFunc func(username, password string) *User) *User {
	value := r.Header.Get("Authorization")
	if len(value) == 0 {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	valueParts := strings.SplitN(value, " ", 2)
	if len(valueParts) != 2 || valueParts[0] != "Basic" {
		http.Error(w, "Invalid authorization header", http.StatusBadRequest)
		return nil
	}

	payload, err := base64.StdEncoding.DecodeString(valueParts[1])
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusBadRequest)
		return nil
	}

	payloadParts := strings.SplitN(string(payload), ":", 2)
	if len(payloadParts) != 2 {
		http.Error(w, "Invalid authorization header", http.StatusBadRequest)
		return nil
	}

	username, password := payloadParts[0], payloadParts[1]
	user := authFunc(username, password)
	if user == nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	return user
}
