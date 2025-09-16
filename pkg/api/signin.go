package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		http.Error(w, `{"error": "Authentication not required"}`, http.StatusBadRequest)
		return
	}

	var auth struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		http.Error(w, `{"error": "Invalid data format"}`, http.StatusBadRequest)
		return
	}

	if auth.Password != password {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid password"})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["pass_hash"] = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	claims["exp"] = time.Now().Add(8 * time.Hour).Unix()
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		http.Error(w, `{"error": "Token generation error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(password), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			storedHash := claims["pass_hash"].(string)
			currentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
			if storedHash != currentHash {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
