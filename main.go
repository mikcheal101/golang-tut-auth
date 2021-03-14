package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

func main() {

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
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func loginEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[*] login invoked!")
}

func registerEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[*] sign up invoked!")
}

func profileEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[*] protected route invoked!")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return next
}
