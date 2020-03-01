package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func getToken() (string, error) {
	userToken := os.Getenv("USER_TOKEN")
	if userToken == "" {
		return "", fmt.Errorf("Missing USER_TOKEN env")
	}
	return userToken, nil
}

func validateToken(encodedAuthToken string) error {
	authTokenBytes, err := base64.StdEncoding.DecodeString(encodedAuthToken)
	authToken := string(authTokenBytes)
	if err != nil {
		return err
	}
	userToken, err := getToken()
	if err != nil {
		return err
	}
	if userToken != authToken {
		return fmt.Errorf("Invalid Token")
	}
	return nil
}

func AuthorizationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if len(authToken) == 0 {
			http.Error(w, "Unauthorized Access", http.StatusUnauthorized)
			return
		}
		authToken = strings.Replace(authToken, "Basic ", "", 1)
		err := validateToken(authToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
