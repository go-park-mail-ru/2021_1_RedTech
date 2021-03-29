package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
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
		http.Error(w, jsonerrors.JSONMessage("json decode"), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userId64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, jsonerrors.JSONMessage("params"), http.StatusBadRequest)
		return
	}
	userId := uint(userId64)
	userUpdate.ID = userId

	currentId, err := getCurrentId(r)
	if err != nil {
		log.Printf("Error while getting current user")
		http.Error(w, jsonerrors.JSONMessage("session"), user.CodeFromError(err))
		return
	} else if currentId != userId {
		log.Printf("Trying to update another user")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), http.StatusBadRequest)
		return
	}

	if err := handler.UUsecase.Update(userUpdate); err != nil {
		log.Printf("Error while updating user")
		http.Error(w, jsonerrors.JSONMessage("invalid update"), http.StatusBadRequest)
		return
	}

	userToSend, err := handler.UUsecase.GetById(userId)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		http.Error(w, jsonerrors.JSONMessage("database"), user.CodeFromError(err))
		return
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		http.Error(w, jsonerrors.JSONMessage("json encode"), http.StatusInternalServerError)
		return
	}
}
