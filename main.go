package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

func handleError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	pgUrl, err := pq.ParseURL("postgres://olretafs:tHhnTvI91wYik-weIGGHTxSqArqpRfTZ@ziggy.db.elephantsql.com:5432/olretafs")
	handleError(err)

	// open pg connection
	db, err = sql.Open("postgres", pgUrl)
	handleError(err)

	router := mux.NewRouter()

	// login route
	router.HandleFunc("/login", loginEndpoint).Methods("POST")

	// signup route
	router.HandleFunc("/signup", registerEndpoint).Methods("POST")

	// profile page - protected route
	router.HandleFunc("/profile", TokenVerifyMiddleware(profileEndpoint)).Methods("GET")

	http.Handle("/", router)

	log.Println("Listening on port: 8000")
	err = http.ListenAndServe(":9090", router)
	handleError(err)
}

func GenerateToken(user User) (string, error) {
	var err error
	secret := "secret" // for tutorial purposes
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
	var user User
	var error Error
	var jwt JWT

	// decode the entered data into user
	json.NewDecoder(req.Body).Decode(&user)

	// validate entry
	if strings.Trim(user.Username, " ") == "" {
		error.Message = "Username is required!"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if strings.Trim(user.Password, " ") == "" {
		error.Message = "Password is required!"
		respondWithError(w, http.StatusBadRequest, error)
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
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(user.Password))
	if err != nil {
		error.Message = "User does not exist!"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	token, err := GenerateToken(user)
	handleError(err)

	jwt.Token = token
	respondWithJson(w, jwt)
}

func respondWithError(w http.ResponseWriter, status int, error Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func respondWithJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func registerEndpoint(w http.ResponseWriter, req *http.Request) {
	var user User
	var error Error

	json.NewDecoder(req.Body).Decode(&user)

	if strings.Trim(user.Username, " ") == "" {
		error.Message = "Username is required!"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if strings.Trim(user.Password, " ") == "" {
		error.Message = "Password is required!"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	handleError(err)
	user.Password = string(hash)

	stmt := "insert into users (username, password) values ($1, $2) RETURNING id;"
	err = db.QueryRow(stmt, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		error.Message = "Invalid user credentials!"
		respondWithError(w, http.StatusInternalServerError, error)
		return
	}

	// reset password to disable the returning of the password to the user
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func profileEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[*] protected route invoked!")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var errorObject Error
		authHeader := r.Header.Get("authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("%v", "There was an error")
				}

				return []byte("secret"), nil
			})

			if err != nil {
				errorObject.Message = err.Error()
				respondWithError(rw, http.StatusUnauthorized, errorObject)
				return
			}

			if token.Valid {
				next.ServeHTTP(rw, r)
			} else {
				errorObject.Message = err.Error()
				respondWithError(rw, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token"
			respondWithError(rw, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
