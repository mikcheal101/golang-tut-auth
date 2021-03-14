package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mikcheal101/golang-tut-auth/models"
	"github.com/mikcheal101/golang-tut-auth/repository"
	"github.com/mikcheal101/golang-tut-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct{}

// method to authenticate a user
func (controller Controller) LoginEndpoint(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var user models.User
		var error models.Error
		var jwt models.JWT
		var userRepo repository.UserRepository

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

		pwd, err := userRepo.AuthUser(db, user)
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
	return func(w http.ResponseWriter, req *http.Request) {
		var user models.User
		var error models.Error
		var userRepo repository.UserRepository

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

		err := userRepo.CreateUser(db, &user)
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
