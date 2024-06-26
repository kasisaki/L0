package server

import (
	mod "L0/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func HandleError(w http.ResponseWriter, status int, err error) {
	errResponse := mod.ErrorResponse{Error: err.Error()}
	responseJSON, marshalErr := json.Marshal(errResponse)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(responseJSON)
}

func HandleGetError(w http.ResponseWriter, err error) bool {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			HandleError(w, 404, errors.New("Запись не найдена"))
			return true
		}
		HandleError(w, http.StatusInternalServerError, err)
		return true
	}
	return false
}

func HandleNormalResponse(w http.ResponseWriter, strct any) {
	responseJSON, err := json.Marshal(strct)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
