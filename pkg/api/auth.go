package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signInRequest struct {
	Password string `json:"password"`
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		respondError(w, http.StatusBadRequest, "authentication not configured")
		return
	}
	if req.Password != pass {
		respondError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	token, err := createToken(pass)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   8 * 3600,
		SameSite: http.SameSiteLaxMode,
	})
	respondJSON(w, map[string]any{"token": token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			respondError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		if !validateToken(cookie.Value, pass) {
			respondError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		next(w, r)
	})
}

func createToken(password string) (string, error) {
	claims := jwt.MapClaims{
		"pwd_hash": passwordHash(password),
		"exp":      time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(password))
}

func validateToken(tokenString, password string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(password), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	hash, ok := claims["pwd_hash"].(string)
	if !ok {
		return false
	}

	return hash == passwordHash(password)
}

func passwordHash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}
