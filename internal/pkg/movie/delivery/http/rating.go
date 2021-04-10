package http

import (
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (handler *MovieHandler) Like(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Log.Warn(fmt.Sprint("Error while getting movie id: ", vars["id"]))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}

	movieId := uint(urlID)

	sess, err := session.GetSession(r)
	if err != nil || session.Manager.Check(sess) != nil {
		log.Log.Warn("Trying to like while unauthorized")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(user.UnauthorizedError))
		return
	}

	err = handler.MUCase.Like(sess.UserID, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while like: %v", err))
		http.Error(w, jsonerrors.JSONMessage("can't set like"), user.CodeFromError(err))
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "success"}); err != nil {
		log.Log.Warn(fmt.Sprintf("Error while encoding JSON: %s", err))
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}

func (handler *MovieHandler) Dislike(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Log.Warn(fmt.Sprint("Error while getting movie id: ", vars["id"]))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}

	movieId := uint(urlID)

	sess, err := session.GetSession(r)
	if err != nil || session.Manager.Check(sess) != nil {
		log.Log.Warn("Trying to dislike while unauthorized")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(user.UnauthorizedError))
		return
	}

	err = handler.MUCase.Dislike(sess.UserID, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while dislike: %v", err))
		http.Error(w, jsonerrors.JSONMessage("can't set dislike"), user.CodeFromError(err))
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "success"}); err != nil {
		log.Log.Warn(fmt.Sprintf("Error while encoding JSON: %s", err))
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
