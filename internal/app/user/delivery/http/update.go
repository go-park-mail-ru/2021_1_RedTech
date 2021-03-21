package http

import (
	"Redioteka/internal/app/domain"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (handler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userUpdate := &domain.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(userUpdate); err != nil {
		log.Printf("Error while unmarshalling JSON")
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userId64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, `{"error":"bad id"}`, http.StatusBadRequest)
		return
	}
	userId := uint(userId64)
	userUpdate.ID = userId

	currentId, err := getCurrentId(r)
	if err != nil {
		log.Printf("Error while getting current user")
		http.Error(w, `{"error":"error while updating user"}`, http.StatusBadRequest)
		return
	} else if currentId != userId {
		log.Printf("Trying to update another user")
		http.Error(w, `{"error":"error while updating user"}`, http.StatusBadRequest)
		return
	}

	if err := handler.UHandler.Update(userUpdate); err != nil {
		log.Printf("Error while updating user")
		http.Error(w, `{"error":"error while updating user"}`, http.StatusBadRequest)
		return
	}

	userToSend, err := handler.UHandler.GetById(userId)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
