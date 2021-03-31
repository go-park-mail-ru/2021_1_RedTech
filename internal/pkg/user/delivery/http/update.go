package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	sess, err := getSession(r)
	if err != nil {
		log.Printf("Error while getting current user session")
		http.Error(w, `{"error":"error while updating user"}`, http.StatusBadRequest)
		return
	} else if session.Manager.Check(sess) != nil || sess.UserID != userId {
		log.Printf("Error while updating user %d", userId)
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
