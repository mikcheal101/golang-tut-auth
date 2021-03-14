package repository

import (
	"database/sql"

	"github.com/mikcheal101/golang-tut-auth/utils"
	"github.com/mikcheal101/golang-tut-auth/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {}

func (repo UserRepository) CreateUser(db *sql.DB, user *models.User) (error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	utils.HandleError(err)
	user.Password = string(hash)

	stmt := "insert into users (username, password) values ($1, $2) RETURNING id;"
	err = db.QueryRow(stmt, user.Username, user.Password).Scan(&user.ID)
	return err
}

func (repo UserRepository) AuthUser(db *sql.DB, user models.User) (string, error) {
	var pwd string
	stmt := "select id, username, password from users where username=$1"
	err := db.QueryRow(stmt, user.Username).Scan(&user.ID, &user.Username, &pwd)
	return pwd, err
}