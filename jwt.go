package main

import (
	"fmt"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Unauthorized"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s AccountStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		userId, err := getId(r)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			return
		}

		account, err := s.GetAccountByID(userId)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			permissionDenied(w)
			return
		}

		accountNumber := int64(claims["accountNumber"].(float64))
		if account.Number != accountNumber {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func createJWT(account *Account) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []
		return []byte(secret), nil
	})
}
