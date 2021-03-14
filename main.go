package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"golang.org/x/crypto/bcrypt"

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

func loginEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[*] login invoked!")
}

func respondWithError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
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
	fmt.Println(hash)
}

func profileEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[*] protected route invoked!")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return next
}
