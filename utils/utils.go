package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/mikcheal101/golang-tut-auth/models"
)

func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func RespondWithJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func HandleError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func GenerateToken(user models.User) (string, error) {
	var err error
	secret := os.Getenv("APP_SECRET")
	// jwt contains header.payload.secret

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"iss":      "course",
	})

	tokenString, err := token.SignedString([]byte(secret))
	HandleError(err)

	return tokenString, nil
}