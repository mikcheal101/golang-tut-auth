package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/mikcheal101/golang-tut-auth/models"
	"github.com/mikcheal101/golang-tut-auth/utils"
)

type Middleware struct {}

// method to verify token
func (middleware Middleware) TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		authHeader := r.Header.Get("authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("%v", "There was an error")
				}

				return []byte(os.Getenv("APP_SECRET")), nil
			})

			if err != nil {
				errorObject.Message = err.Error()
				utils.RespondWithError(rw, http.StatusUnauthorized, errorObject)
				return
			}

			if token.Valid {
				next.ServeHTTP(rw, r)
			} else {
				errorObject.Message = err.Error()
				utils.RespondWithError(rw, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token"
			utils.RespondWithError(rw, http.StatusUnauthorized, errorObject)
			return
		}
	})
}