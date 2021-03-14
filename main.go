package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikcheal101/golang-tut-auth/controllers"
	"github.com/mikcheal101/golang-tut-auth/driver"
	"github.com/mikcheal101/golang-tut-auth/middleware"
	"github.com/mikcheal101/golang-tut-auth/utils"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	db = driver.ConnectDB()

	router := mux.NewRouter()
	controller := controllers.Controller{}
	middleware := middleware.Middleware{}

	// login route
	router.HandleFunc("/login", controller.LoginEndpoint(db)).Methods("POST")

	// signup route
	router.HandleFunc("/signup", controller.RegisterEndpoint(db)).Methods("POST")

	// profile page - protected route
	router.HandleFunc("/profile", middleware.TokenVerifyMiddleware(controller.ProfileEndpoint(db))).Methods("GET")

	http.Handle("/", router)

	log.Println("Listening on port: 8000")
	err := http.ListenAndServe(":9090", router)
	utils.HandleError(err)
}
