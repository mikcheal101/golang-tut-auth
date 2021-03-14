package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mikcheal101/golang-tut-auth/models"
	"github.com/mikcheal101/golang-tut-auth/utils"
	"golang.org/x/crypto/bcrypt"
)


type Controller struct {}

// method to authenticate a user
func (controller Controller) LoginEndpoint(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, req *http.Request) {
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
		
		token, err := utils.GenerateToken(user)
		utils.HandleError(err)

		user.Password = ""
		w.Header().Set("Authorization", token)
		utils.HandleError(err)
	
		jwt.Token = token
		utils.RespondWithJson(w, jwt)
	}
}

// method to register a user
func (controller Controller) RegisterEndpoint(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, req *http.Request) {
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
		utils.HandleError(err)
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
		utils.RespondWithJson(w, user)
	}
}



