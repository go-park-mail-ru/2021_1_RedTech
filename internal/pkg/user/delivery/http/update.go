package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (handler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userUpdate := &domain.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(userUpdate); err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot parse update form: %s", err))
		http.Error(w, jsonerrors.JSONDecode, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userId64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot get user id: %s", err))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	userId := uint(userId64)
	userUpdate.ID = userId

	sess, err := session.GetSession(r)
	if err != nil {
		log.Log.Warn("Error while getting current user session")
		http.Error(w, jsonerrors.Session, user.CodeFromError(err))
		return
	} else if session.Manager.Check(sess) != nil || sess.UserID != userId {
		log.Log.Warn(fmt.Sprintf("Error while updating user %d", userId))
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), http.StatusBadRequest)
		return
	}

	if err := handler.UUsecase.Update(userUpdate); err != nil {
		log.Log.Warn(fmt.Sprintf("Error while updating user %d", userId))
		http.Error(w, jsonerrors.JSONMessage("invalid update"), user.CodeFromError(err))
		return
	}

	userToSend, err := handler.UUsecase.GetById(userId)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		http.Error(w, jsonerrors.JSONMessage("database"), user.CodeFromError(err))
		return
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
