package controllers

import (
	"database/sql"
	"net/http"

	"github.com/mikcheal101/golang-tut-auth/utils"
)

func (controller Controller) ProfileEndpoint(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, req *http.Request) {
		utils.RespondWithJson(w, "it works")
	}
}
