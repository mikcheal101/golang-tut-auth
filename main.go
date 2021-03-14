package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mikcheal101/golang-tut-auth/driver"
	"github.com/mikcheal101/golang-tut-auth/models"
	"github.com/mikcheal101/golang-tut-auth/utils"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func handleError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func init() {
	gotenv.Load()
}

func main() {
	db = driver.ConnectDB()

	router := mux.NewRouter()

	// login route
	router.HandleFunc("/login", loginEndpoint).Methods("POST")

	// signup route
	router.HandleFunc("/signup", registerEndpoint).Methods("POST")

	// profile page - protected route
	router.HandleFunc("/profile", TokenVerifyMiddleware(profileEndpoint)).Methods("GET")

	http.Handle("/", router)

	log.Println("Listening on port: 8000")
	err := http.ListenAndServe(":9090", router)
	handleError(err)
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
	handleError(err)

	return tokenString, nil
}

func loginEndpoint(w http.ResponseWriter, req *http.Request) {
	var user models.User
	var error models.Error
	var jwt models.JWT

	// decode the entered data into user
	json.NewDecoder(req.Body).Decode(&user)

	// validate entry
	if strings.Trim(user.Username, " ") == "" {
		error.Message = "Username is required!"
		utils.RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	if strings.Trim(user.Password, " ") == "" {
		error.Message = "Password is required!"
		utils.RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	var pwd string
	stmt := "select id, username, password from users where username=$1"
	err := db.QueryRow(stmt, user.Username).Scan(&user.ID, &user.Username, &pwd)
	if err != nil {
		error.Message = err.Error()
		if err == sql.ErrNoRows {
			error.Message = "User does not exist!"
		}
		utils.RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(user.Password))
	if err != nil {
		error.Message = "User does not exist!"
		utils.RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	token, err := GenerateToken(user)
	handleError(err)

	jwt.Token = token
	utils.RespondWithJson(w, jwt)
}

func registerEndpoint(w http.ResponseWriter, req *http.Request) {
	var user models.User
	var error models.Error

	json.NewDecoder(req.Body).Decode(&user)

	if strings.Trim(user.Username, " ") == "" {
		error.Message = "Username is required!"
		utils.RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	if strings.Trim(user.Password, " ") == "" {
		error.Message = "Password is required!"
		utils.RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	handleError(err)
	user.Password = string(hash)

	stmt := "insert into users (username, password) values ($1, $2) RETURNING id;"
	err = db.QueryRow(stmt, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		error.Message = "User already exists!"
		utils.RespondWithError(w, http.StatusInternalServerError, error)
		return
	}

	// reset password to disable the returning of the password to the user
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func profileEndpoint(w http.ResponseWriter, req *http.Request) {
	utils.RespondWithJson(w, "it works")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
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
