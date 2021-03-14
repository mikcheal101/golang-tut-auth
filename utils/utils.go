package utils

import (
	"encoding/json"
	"net/http"

	"github.com/mikcheal101/golang-tut-auth/models"
)


func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func RespondWithJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
